/*
Copyright © 2024 EdgeFarm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"tideland.dev/go/wait"

	strutil "github.com/openyurtio/openyurt/pkg/util/strings"
	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"

	"github.com/edgefarm/edgefarm/pkg/constants"
	deploy "github.com/edgefarm/edgefarm/pkg/deploy"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/k8s/addons"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	stringsx "github.com/icza/gox/stringsx"
)

func NewCreateCommand(out io.Writer) *cobra.Command {
	o := newKindOptions()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the local edgefarm cluster",
		RunE: func(cmd *cobra.Command, arguments []string) error {
			_, err := state.GetState()
			if err != nil {
				return err
			}
			if err := shared.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if shared.Args.Only.CoreDNS || shared.Args.Only.Flannel || shared.Args.Only.KubeProxy {
				shared.Args.Skip = shared.ConvertOnlyToSkip(shared.Args.Only)
			}
			if err := o.Validate(); err != nil {
				return err
			}
			initializer, err := newKindInitializer(out, o.Config())
			if err != nil {
				return err
			}
			if err := initializer.Run(); err != nil {
				return err
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.SetOut(out)
	shared.AddSharedFlags(cmd.Flags())
	addFlagsForCreate(cmd.Flags(), o)
	deploy.AddFlagsForDeploy(cmd.Flags())
	return cmd
}

func init() {
	localClusterCmd.AddCommand(NewCreateCommand(os.Stdout))
}

type kindOptions struct {
	KindConfigPath    string
	WorkerNodesNum    int
	EdgeNodesNum      int
	ClusterName       string
	CloudNodes        string
	OpenYurtVersion   string
	KubernetesVersion string
	UseLocalImages    bool
	KubeConfig        string
	IgnoreError       bool
	NodeImage         string
}

func newKindOptions() *kindOptions {
	return &kindOptions{
		WorkerNodesNum:    1,
		EdgeNodesNum:      2,
		ClusterName:       "edgefarm",
		OpenYurtVersion:   "v1.4.0",
		KubernetesVersion: "v1.22.7",
		UseLocalImages:    false,
		IgnoreError:       true,
		CloudNodes:        "edgefarm-control-plane,edgefarm-worker",
		NodeImage:         "ghcr.io/edgefarm/edgefarm/kind-node:v1.22.7-systemd",
	}
}

func checkFreePort(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
		conn.Close()
		return false
	}
	return true
}

func (o *kindOptions) Validate() error {
	if runtime.GOOS == "linux" {
		notGood := false
		f, err := os.ReadFile("/proc/sys/fs/inotify/max_user_instances")
		if err != nil {
			return err
		}

		v, err := strconv.Atoi(stringsx.Clean(string(f)))
		if err != nil {
			return err
		}
		if v < 512 {
			notGood = true
			klog.Errorln("the value of /proc/sys/fs/inotify/max_user_instances must be greater or equal than 512")
		}

		f, err = os.ReadFile("/proc/sys/fs/inotify/max_user_watches")
		if err != nil {
			return err
		}

		v, err = strconv.Atoi(stringsx.Clean(string(f)))
		if err != nil {
			return err
		}

		if v < 524288 {
			notGood = true
			klog.Errorln("the value of /proc/sys/fs/inotify/max_user_watches must be greater or equal than 524288")
		}

		if notGood {
			return errors.New("follow https://kind.sigs.k8s.io/docs/user/known-issues/#pod-errors-due-to-too-many-open-files to fix this issue")
		}
	}

	if o.WorkerNodesNum < 1 {
		return fmt.Errorf("the number of nodes must be greater than 0")
	}
	if !checkFreePort(shared.Ports.HostApiServerPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostApiServerPort)
	}
	if !checkFreePort(shared.Ports.HostNatsPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostNatsPort)
	}
	if !checkFreePort(shared.Ports.HostHttpPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostHttpPort)
	}
	if !checkFreePort(shared.Ports.HostHttpsPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostHttpsPort)
	}
	return nil
}

// Config should be called after Validate
// It will generate a config for Initializer
func (o *kindOptions) Config() *initializerConfig {
	controlPlaneNode, workerNodes := getNodeNamesOfKindCluster(o.ClusterName, o.WorkerNodesNum, o.EdgeNodesNum)
	allNodes := append(workerNodes, controlPlaneNode)

	// prepare kindConfig.CloudNodes and kindConfig.EdgeNodes
	cloudNodes := sets.NewString()
	if o.CloudNodes != "" {
		for _, node := range strings.Split(o.CloudNodes, ",") {
			if !strutil.IsInStringLst(allNodes, node) {
				klog.Fatalf("node %s will not be in the cluster", node)
			}
			cloudNodes = cloudNodes.Insert(node)
		}
	}
	// any node not be specified as cloud node will be recognized as edge node
	edgeNodes := sets.NewString()
	for _, node := range allNodes {
		if !cloudNodes.Has(node) {
			edgeNodes = edgeNodes.Insert(node)
		}
	}

	return &initializerConfig{
		CloudNodes:        cloudNodes.List(),
		EdgeNodes:         edgeNodes.List(),
		KindConfigPath:    o.KindConfigPath,
		WorkerNodesNum:    o.WorkerNodesNum,
		EdgeNodesNum:      o.EdgeNodesNum,
		NodeImage:         o.NodeImage,
		ClusterName:       o.ClusterName,
		KubernetesVersion: o.KubernetesVersion,
		UseLocalImage:     o.UseLocalImages,
	}
}

func addFlagsForCreate(flagset *pflag.FlagSet, o *kindOptions) {
	flagset.IntVar(&o.EdgeNodesNum, "edge-node-num", o.EdgeNodesNum, "Specify the edge node number of the kind cluster.")
	flagset.IntVar(&shared.Ports.HostApiServerPort, "host-api-server-port", shared.Ports.HostApiServerPort, "Specify the port of host api server.")
	flagset.IntVar(&shared.Ports.HostNatsPort, "host-nats-port", shared.Ports.HostNatsPort, "Specify the port of nats to be mapped to.")
	flagset.IntVar(&shared.Ports.HostHttpPort, "host-http-port", shared.Ports.HostHttpPort, "Specify the port of http server to be mapped to.")
	flagset.IntVar(&shared.Ports.HostHttpsPort, "host-https-port", shared.Ports.HostHttpsPort, "Specify the port of https server to be mapped to.")
	flagset.BoolVar(&shared.Args.Deploy, "deploy", false, "Deploy the cluster after creation.")
	flagset.BoolVar(&shared.Args.Skip.CoreDNS, "skip-coredns", false, "Skip installing CoreDNS. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Skip.KubeProxy, "skip-kube-proxy", false, "Skip installing kube-proxy. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Skip.Flannel, "skip-flannel", false, "Skip installing flannel. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.CoreDNS, "only-coredns", false, "Only install CoreDNS. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.Flannel, "only-flannel", false, "Only install flannel. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.KubeProxy, "only-kube-proxy", false, "Only install kube-proxy. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
}

type initializerConfig struct {
	CloudNodes        []string
	EdgeNodes         []string
	KindConfigPath    string
	WorkerNodesNum    int
	EdgeNodesNum      int
	ClusterName       string
	KubernetesVersion string
	NodeImage         string
	UseLocalImage     bool
}

type Initializer struct {
	initializerConfig
	out      io.Writer
	operator *kindoperator.KindOperator
}

func newKindInitializer(out io.Writer, cfg *initializerConfig) (*Initializer, error) {
	k, err := kindoperator.NewKindOperator(shared.KubeConfig)
	if err != nil {
		return nil, err
	}
	return &Initializer{
		initializerConfig: *cfg,
		out:               out,
		operator:          k,
	}, nil
}

func (ki *Initializer) Run() error {
	klog.Info("Start to prepare config for kind")
	config, err := ki.prepareKindConfigFile()
	if err != nil {
		return err
	}

	klog.Info("Start to create cluster with kind")
	if err := ki.operator.KindCreateClusterWithConfig(config); err != nil {
		return err
	}

	klog.Info("Cluster created")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", shared.KubeConfig)
	if err != nil {
		return err
	}
	shared.KubeConfigRestConfig = kubeconfig

	if err := k8s.PrepareEdgeNodes(); err != nil {
		return err
	}

	if !shared.Args.Skip.CoreDNS {
		klog.Infof("Deploy cluster coredns packages")
		if err := addons.ReplaceCoreDNS(); err != nil {
			return err
		}
	}
	if !shared.Args.Skip.Flannel {
		klog.Infof("Deploy cluster flannel packages")
		if err := packages.Install(packages.Flannel); err != nil {
			return err
		}
		if err := WaitForBootstrapConditions(time.Minute * 5); err != nil {
			return err
		}
	}
	if !shared.Args.Skip.KubeProxy {
		klog.Infof("Deploy cluster kube-proxy packages")
		if err := addons.ReplaceKubeProxy(); err != nil {
			return err
		}
	}

	if shared.Args.Deploy {
		if err := deploy.Deploy(); err != nil {
			return err
		}
	} else {
		green := color.New(color.FgHiGreen)
		yellow := color.New(color.FgHiYellow)
		green.Printf("The local cluster has been created.\nRun ")
		yellow.Printf("  $ local-up deploy")
		green.Printf(" to deploy EdgeFarm components and its dependencies.\nHave a look at the arguments using '--help'.")
	}
	return nil
}

// func (ki *Initializer) prepareImages() error {
// 	if !ki.UseLocalImage {
// 		return nil
// 	}
// 	// load images of cloud components to cloud nodes
// 	if err := ki.loadImagesToKindNodes([]string{
// 		ki.YurtHubImage,
// 		ki.YurtManagerImage,
// 		ki.NodeServantImage,
// 	}, ki.CloudNodes); err != nil {
// 		return err
// 	}

// 	// load images of edge components to edge nodes
// 	if err := ki.loadImagesToKindNodes([]string{
// 		ki.YurtHubImage,
// 		ki.NodeServantImage,
// 	}, ki.EdgeNodes); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ki *Initializer) prepareKindNodeImage() error {
// 	kindVer, err := ki.operator.KindVersion()
// 	if err != nil {
// 		return err
// 	}
// 	ki.NodeImage = kindNodeImageMap[kindVer][ki.KubernetesVersion]
// 	if len(ki.NodeImage) == 0 {
// 		return fmt.Errorf("failed to get node image by kind version= %s and kubernetes version= %s", kindVer, ki.KubernetesVersion)
// 	}

// 	return nil
// }

func (ki *Initializer) prepareKindConfigFile() ([]byte, error) {
	kindConfigContent, err := tmplutil.SubsituteTemplate(constants.KindConfigTemplate, map[string]string{
		"kubernetes_version":   ki.KubernetesVersion,
		"kind_node_image":      ki.NodeImage,
		"cluster_name":         ki.ClusterName,
		"host_api_server_port": fmt.Sprintf("%d", shared.Ports.HostApiServerPort),
	})
	if err != nil {
		return nil, err
	}

	// add additional worker entries into kind config file according to NodesNum
	for num := 0; num < ki.WorkerNodesNum; num++ {
		worker, err := tmplutil.SubsituteTemplate(constants.KindWorkerRoleTemplate, map[string]string{
			"kind_node_image": ki.NodeImage,
			"host_nats_port":  fmt.Sprintf("%d", shared.Ports.HostNatsPort),
			"host_http_port":  fmt.Sprintf("%d", shared.Ports.HostHttpPort),
			"host_https_port": fmt.Sprintf("%d", shared.Ports.HostHttpsPort),
		})
		if err != nil {
			return nil, err
		}
		kindConfigContent = strings.Join([]string{kindConfigContent, worker}, "\n")
	}

	for num := 0; num < ki.EdgeNodesNum; num++ {
		worker, err := tmplutil.SubsituteTemplate(constants.KindEdgeRole, map[string]string{
			"kind_node_image": ki.NodeImage,
		})
		if err != nil {
			return nil, err
		}
		kindConfigContent = strings.Join([]string{kindConfigContent, worker}, "\n")
	}
	return []byte(kindConfigContent), nil
}

// func (ki *Initializer) deployOpenYurt() error {
// 	// dir, err := os.Getwd()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// converter := &yurtinit.DeployOpenYurt{
// 	// 	// RootDir:           dir,
// 	// 	ComponentsBuilder: kubeutil.NewBuilder(shared.KubeConfig),
// 	// 	ClientSet:         ki.kubeClient,
// 	// 	// CloudNodes:                ki.CloudNodes,
// 	// 	// EdgeNodes:                 ki.EdgeNodes,
// 	// 	// WaitServantJobTimeout:     kubeutil.DefaultWaitServantJobTimeout,
// 	// 	YurthubHealthCheckTimeout: 2 * time.Minute,
// 	// 	KubeConfigPath:            shared.KubeConfig,
// 	// 	YurthubImage:              fmt.Sprintf(constants.YurtHubImageFormat, constants.OpenYurtVersion),
// 	// 	YurtManagerImage:          fmt.Sprintf(constants.YurtManagerImageFormat, constants.OpenYurtVersion),
// 	// 	NodeServantImage:          fmt.Sprintf(constants.NodeServantImageFormat, constants.OpenYurtVersion),
// 	// 	EnableDummyIf:             ki.EnableDummyIf,
// 	// }
// 	// if err := converter.Run(); err != nil {
// 	// 	klog.Errorf("errors occurred when deploying openyurt components")
// 	// 	return err
// 	// }
// 	return nil
// }

// func (ki *Initializer) loadImagesToKindNodes(images, nodes []string) error {
// 	for _, image := range images {
// 		if image == "" {
// 			// if image == "", it's the responsibility of kind to pull images from registry.
// 			continue
// 		}
// 		if err := ki.operator.KindLoadDockerImage(ki.out, ki.ClusterName, image, nodes); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// getNodeNamesOfKindCluster will generate all nodes will be in the kind cluster.
// It depends on the naming machanism of kind:
// one control-plane node: ${clusterName}-control-plane
// serval worker nodes: ${clusterName}-worker, ${clusterName}-worker2, ${clusterName}-worker3...
func getNodeNamesOfKindCluster(clusterName string, workerNodesNum int, edgeNodesNum int) (string, []string) {
	controlPlaneNode := fmt.Sprintf("%s-control-plane", clusterName)
	workerNodes := []string{}
	if workerNodesNum >= 1 {
		workerNodes = append(workerNodes, strings.Join([]string{clusterName, "worker"}, "-"))
	}
	for cnt := 0; cnt < (workerNodesNum-1)+edgeNodesNum; cnt++ {
		workerNodes = append(workerNodes, fmt.Sprintf("%s-worker%d", clusterName, 2+cnt))
	}
	return controlPlaneNode, workerNodes
}

func WaitForBootstrapConditions(stepTimeout time.Duration) error {
	ticker := wait.MakeExpiringIntervalTicker(time.Second, stepTimeout)

	// Checks for ready state of all nodes
	nodesCondition := func() (bool, error) {
		nodes, err := k8s.GetAllNodes()
		if err != nil {
			return false, err
		}
		for _, node := range nodes {
			for _, condition := range node.Status.Conditions {
				if condition.Type == v1.NodeReady && condition.Status != v1.ConditionTrue {
					return false, nil
				}
			}
		}
		return true, nil
	}
	err := wait.Poll(context.Background(), ticker, nodesCondition)
	if err != nil {
		return err
	}

	return nil
}
