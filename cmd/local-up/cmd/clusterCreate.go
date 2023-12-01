/*
Copyright © 2023 EdgeFarm Authors

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
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/openyurtio/openyurt/pkg/projectinfo"
	strutil "github.com/openyurtio/openyurt/pkg/util/strings"
	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
	yurtconstantes "github.com/openyurtio/openyurt/test/e2e/cmd/init/constants"
	kubeutil "github.com/openyurtio/openyurt/test/e2e/cmd/init/util/kubernetes"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/constants"
	yurtinit "github.com/edgefarm/edgefarm/pkg/init"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/packages"
)

var (
	validKubernetesVersions = []string{
		"v1.17",
		"v1.18",
		"v1.19",
		"v1.20",
		"v1.21",
		"v1.22",
		"v1.23",
	}

	AllValidOpenYurtVersions = append(projectinfo.Get().AllVersions, "latest")

	kindNodeImageMap = map[string]map[string]string{
		"v0.12.0": {
			"v1.22": "ghcr.io/edgefarm/edgefarm/kind-node:v1.22.7",
		},
	}

	yurtHubImageFormat     = "openyurt/yurthub:%s"
	yurtManagerImageFormat = "openyurt/yurt-manager:%s"
	nodeServantImageFormat = "openyurt/node-servant:%s"
)

func NewCreateCommand(out io.Writer) *cobra.Command {
	o := newKindOptions()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the local edgefarm cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			initializer := newKindInitializer(out, o.Config())
			if err := initializer.Run(); err != nil {
				return err
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.SetOut(out)
	addFlags(cmd.Flags(), o)
	return cmd
}

func init() {
	localClusterCmd.AddCommand(NewCreateCommand(os.Stdout))
}

type kindOptions struct {
	KindConfigPath string
	WorkerNodesNum int
	EdgeNodesNum   int

	ClusterName       string
	CloudNodes        string
	OpenYurtVersion   string
	KubernetesVersion string
	UseLocalImages    bool
	KubeConfig        string
	IgnoreError       bool
	EnableDummyIf     bool
	DisableDefaultCNI bool
}

func newKindOptions() *kindOptions {
	return &kindOptions{
		KindConfigPath:    fmt.Sprintf("%s/kindconfig.yaml", yurtconstantes.TmpDownloadDir),
		WorkerNodesNum:    1,
		EdgeNodesNum:      2,
		ClusterName:       "edgefarm",
		OpenYurtVersion:   "v1.3.4",
		KubernetesVersion: "v1.22",
		UseLocalImages:    false,
		IgnoreError:       true,
		EnableDummyIf:     true,
		DisableDefaultCNI: true,
		CloudNodes:        "edgefarm-control-plane,edgefarm-worker",
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
	if o.WorkerNodesNum < 1 {
		return fmt.Errorf("the number of nodes must be greater than 0")
	}
	if err := validateKubernetesVersion(o.KubernetesVersion); err != nil {
		return err
	}
	if err := validateOpenYurtVersion(o.OpenYurtVersion, o.IgnoreError); err != nil {
		return err
	}
	if !checkFreePort(args.Ports.HostApiServerPort) {
		return fmt.Errorf("port %d is already used", args.Ports.HostApiServerPort)
	}
	if !checkFreePort(args.Ports.HostNatsPort) {
		return fmt.Errorf("port %d is already used", args.Ports.HostNatsPort)
	}
	if !checkFreePort(args.Ports.HostHttpPort) {
		return fmt.Errorf("port %d is already used", args.Ports.HostHttpPort)
	}
	if !checkFreePort(args.Ports.HostHttpsPort) {
		return fmt.Errorf("port %d is already used", args.Ports.HostHttpsPort)
	}
	if !checkFreePort(args.Ports.HostVPNPort) {
		return fmt.Errorf("port %d is already used", args.Ports.HostVPNPort)
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
	cloudNodes = cloudNodes.Insert(controlPlaneNode)
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

	// prepare kindConfig.KubeConfig
	kubeConfigPath := o.KubeConfig
	if kubeConfigPath == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
			klog.V(1).Infof("--kube-config is not specified, %s will be used.", kubeConfigPath)
		} else {
			klog.Fatal("failed to get ${HOME} env when using default kubeconfig path")
		}
	}

	return &initializerConfig{
		CloudNodes:     cloudNodes.List(),
		EdgeNodes:      edgeNodes.List(),
		KindConfigPath: o.KindConfigPath,
		KubeConfig:     kubeConfigPath,
		WorkerNodesNum: o.WorkerNodesNum,
		EdgeNodesNum:   o.EdgeNodesNum,

		ClusterName:       o.ClusterName,
		KubernetesVersion: o.KubernetesVersion,
		UseLocalImage:     o.UseLocalImages,
		YurtHubImage:      fmt.Sprintf(yurtHubImageFormat, o.OpenYurtVersion),
		YurtManagerImage:  fmt.Sprintf(yurtManagerImageFormat, o.OpenYurtVersion),
		NodeServantImage:  fmt.Sprintf(nodeServantImageFormat, o.OpenYurtVersion),
		EnableDummyIf:     o.EnableDummyIf,
		DisableDefaultCNI: o.DisableDefaultCNI,
	}
}

var (
	skipApplications        bool
	skipNetwork             bool
	skipMonitor             bool
	skipClusterDependencies bool
	skipBase                bool
)

func addFlags(flagset *pflag.FlagSet, o *kindOptions) {
	flagset.IntVar(&o.EdgeNodesNum, "edge-node-num", o.EdgeNodesNum, "Specify the edge node number of the kind cluster.")
	flagset.StringVar(&o.KubeConfig, "kube-config", o.KubeConfig, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	flagset.IntVar(&args.Ports.HostApiServerPort, "host-api-server-port", args.Ports.HostApiServerPort, "Specify the port of host api server.")
	flagset.IntVar(&args.Ports.HostNatsPort, "host-nats-port", args.Ports.HostNatsPort, "Specify the port of nats to be mapped to.")
	flagset.IntVar(&args.Ports.HostHttpPort, "host-http-port", args.Ports.HostHttpPort, "Specify the port of http server to be mapped to.")
	flagset.IntVar(&args.Ports.HostHttpsPort, "host-https-port", args.Ports.HostHttpsPort, "Specify the port of https server to be mapped to.")
	flagset.IntVar(&args.Ports.HostVPNPort, "host-vpn-port", args.Ports.HostVPNPort, "Specify the port for the local VPN.")
	flagset.StringVar(&args.Interface, "interface", "", "Network interface to connect to physical edge nodes. This is probably the same interface that is used to connect to the internet. If unset, defaults to the first default routes' interface.")
	flagset.BoolVar(&skipApplications, "skip-applications", false, "Skip installing edgefarm.applications.")
	flagset.BoolVar(&skipNetwork, "skip-network", false, "Skip installing edgefarm.network.")
	flagset.BoolVar(&skipMonitor, "skip-monitor", false, "Skip installing edgefarm.monitor.")
	flagset.BoolVar(&skipClusterDependencies, "skip-cluster-dependencies", false, "Skip installing edgefarm.cluster-dependencies.")
	flagset.BoolVar(&skipBase, "skip-base", false, "Skip installing base packages for edgefarm.")
}

type initializerConfig struct {
	CloudNodes        []string
	EdgeNodes         []string
	KindConfigPath    string
	KubeConfig        string
	WorkerNodesNum    int
	EdgeNodesNum      int
	ClusterName       string
	KubernetesVersion string
	NodeImage         string
	UseLocalImage     bool
	YurtHubImage      string
	YurtManagerImage  string
	NodeServantImage  string
	EnableDummyIf     bool
	DisableDefaultCNI bool
}

type Initializer struct {
	initializerConfig
	out        io.Writer
	operator   *kindoperator.KindOperator
	kubeClient kubeclientset.Interface
}

func newKindInitializer(out io.Writer, cfg *initializerConfig) *Initializer {
	return &Initializer{
		initializerConfig: *cfg,
		out:               out,
		operator:          kindoperator.NewKindOperator("", cfg.KubeConfig),
	}
}

func (ki *Initializer) Run() error {
	klog.Info("Check kind")
	if err := ki.operator.GetKindPath(); err != nil {
		return err
	}

	klog.Info("Start to prepare kind node image")
	if err := ki.prepareKindNodeImage(); err != nil {
		return err
	}

	klog.Info("Start to prepare config file for kind")
	if err := ki.prepareKindConfigFile(ki.KindConfigPath); err != nil {
		return err
	}

	klog.Info("Start to create cluster with kind")
	if err := ki.operator.KindCreateClusterWithConfig(ki.out, ki.KindConfigPath); err != nil {
		return err
	}

	klog.Info("Start to prepare kube client")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", ki.KubeConfig)
	if err != nil {
		return err
	}
	ki.kubeClient, err = kubeclientset.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}
	klog.Infof("Deploy cluster bootstrap stage 1 packages")
	if err := packages.InstallAndWaitBootstrapStage1(); err != nil {
		return err
	}

	klog.Info("Start to prepare OpenYurt images for kind cluster")
	if err := ki.prepareImages(); err != nil {
		return err
	}

	klog.Info("Start to deploy OpenYurt components")
	if err := ki.deployOpenYurt(); err != nil {
		return err
	}

	klog.Infof("Start to configure coredns to adapt OpenYurt")
	if err := ki.configureAddons(); err != nil {
		return err
	}

	klog.Infof("Prepare edge nodes")
	if err := k8s.PrepareEdgeNodes(); err != nil {
		return err
	}

	klog.Infof("Deploy cluster bootstrap stage 2 packages")
	if err := packages.Install(packages.ClusterBootstrapStage2); err != nil {
		return err
	}

	if !skipClusterDependencies {
		klog.Infof("Deploy cluster dependencies")
		if err := packages.Install(packages.ClusterDependencies); err != nil {
			return err
		}
	}

	if !skipBase {
		klog.Infof("Deploy edgefarm base packages")
		if err := packages.Install(packages.Base); err != nil {
			return err
		}
	}

	if !skipNetwork {
		klog.Infof("Deploy edgefarm network packages")
		if err := packages.Install(packages.Network); err != nil {
			return err
		}
	}

	if !skipApplications {
		klog.Infof("Deploy edgefarm applications packages")
		if err := packages.Install(packages.Applications); err != nil {
			return err
		}
	}

	if !skipMonitor {
		klog.Infof("Deploy edgefarm monitor packages")
		if err := packages.Install(packages.Monitor); err != nil {
			return err
		}
	}

	return nil
}

func (ki *Initializer) prepareImages() error {
	if !ki.UseLocalImage {
		return nil
	}
	// load images of cloud components to cloud nodes
	if err := ki.loadImagesToKindNodes([]string{
		ki.YurtHubImage,
		ki.YurtManagerImage,
		ki.NodeServantImage,
	}, ki.CloudNodes); err != nil {
		return err
	}

	// load images of edge components to edge nodes
	if err := ki.loadImagesToKindNodes([]string{
		ki.YurtHubImage,
		ki.NodeServantImage,
	}, ki.EdgeNodes); err != nil {
		return err
	}

	return nil
}

func (ki *Initializer) prepareKindNodeImage() error {
	kindVer, err := ki.operator.KindVersion()
	if err != nil {
		return err
	}
	ki.NodeImage = kindNodeImageMap[kindVer][ki.KubernetesVersion]
	if len(ki.NodeImage) == 0 {
		return fmt.Errorf("failed to get node image by kind version= %s and kubernetes version= %s", kindVer, ki.KubernetesVersion)
	}

	return nil
}

func (ki *Initializer) prepareKindConfigFile(kindConfigPath string) error {
	kindConfigDir := filepath.Dir(kindConfigPath)
	if err := os.MkdirAll(kindConfigDir, yurtconstantes.DirMode); err != nil {
		return err
	}
	kindConfigContent, err := tmplutil.SubsituteTemplate(constants.OpenYurtKindConfig, map[string]string{
		"kind_node_image":      ki.NodeImage,
		"cluster_name":         ki.ClusterName,
		"disable_default_cni":  fmt.Sprintf("%v", ki.DisableDefaultCNI),
		"host_api_server_port": fmt.Sprintf("%d", args.Ports.HostApiServerPort),
		"host_vpn_port":        fmt.Sprintf("%d", args.Ports.HostVPNPort),
	})
	if err != nil {
		return err
	}

	// add additional worker entries into kind config file according to NodesNum
	for num := 0; num < ki.WorkerNodesNum; num++ {
		worker, err := tmplutil.SubsituteTemplate(constants.KindWorkerRole, map[string]string{
			"kind_node_image": ki.NodeImage,
			"host_nats_port":  fmt.Sprintf("%d", args.Ports.HostNatsPort),
			"host_http_port":  fmt.Sprintf("%d", args.Ports.HostHttpPort),
			"host_https_port": fmt.Sprintf("%d", args.Ports.HostHttpsPort),
		})
		if err != nil {
			return err
		}
		kindConfigContent = strings.Join([]string{kindConfigContent, worker}, "\n")
	}

	for num := 0; num < ki.EdgeNodesNum; num++ {
		worker, err := tmplutil.SubsituteTemplate(constants.KindEdgeRole, map[string]string{
			"kind_node_image": ki.NodeImage,
		})
		if err != nil {
			return err
		}
		kindConfigContent = strings.Join([]string{kindConfigContent, worker}, "\n")
	}

	if err = os.WriteFile(kindConfigPath, []byte(kindConfigContent), yurtconstantes.FileMode); err != nil {
		return err
	}
	klog.V(1).Infof("generated new kind config file at %s", kindConfigPath)
	return nil
}

func (ki *Initializer) configureAddons() error {
	if err := ki.configureCoreDnsAddon(); err != nil {
		return err
	}

	// re-construct kube-proxy pods
	podList, err := ki.kubeClient.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for i := range podList.Items {
		switch {
		case strings.HasPrefix(podList.Items[i].Name, "kube-proxy"):
			// delete pod
			propagation := metav1.DeletePropagationForeground
			err = ki.kubeClient.CoreV1().Pods("kube-system").Delete(context.TODO(), podList.Items[i].Name, metav1.DeleteOptions{
				PropagationPolicy: &propagation,
			})
			if err != nil {
				klog.Errorf("failed to delete pod(%s), %v", podList.Items[i].Name, err)
			}
		default:
		}
	}

	// If we disable default cni, nodes will not be ready and the coredns pod always be in pending.
	// The health check for coreDNS should be done by someone who will install CNI.
	if !ki.DisableDefaultCNI {
		// wait for coredns pods available
		for {
			select {
			case <-time.After(10 * time.Second):
				dnsDp, err := ki.kubeClient.AppsV1().Deployments("kube-system").Get(context.TODO(), "coredns", metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("failed to get coredns deployment when waiting for available, %v", err)
				}

				if dnsDp.Status.ObservedGeneration < dnsDp.Generation {
					klog.Infof("waiting for coredns generation(%d) to be observed. now observed generation is %d", dnsDp.Generation, dnsDp.Status.ObservedGeneration)
					continue
				}

				if *dnsDp.Spec.Replicas != dnsDp.Status.AvailableReplicas {
					klog.Infof("waiting for coredns replicas(%d) to be ready, now %d pods available", *dnsDp.Spec.Replicas, dnsDp.Status.AvailableReplicas)
					continue
				}
				klog.Info("coredns deployment configuration is completed")
				return nil
			}
		}
	}
	return nil
}

func (ki *Initializer) configureCoreDnsAddon() error {
	// configure hostname service topology for kube-dns service
	svc, err := ki.kubeClient.CoreV1().Services("kube-system").Get(context.TODO(), "kube-dns", metav1.GetOptions{})
	if err != nil {
		return err
	}

	topologyChanged := false
	if svc != nil {
		if svc.Annotations == nil {
			svc.Annotations = make(map[string]string)
		}

		if val, ok := svc.Annotations["openyurt.io/topologyKeys"]; ok && val == "kubernetes.io/hostname" {
			// topology annotation does not need to change
		} else {
			svc.Annotations["openyurt.io/topologyKeys"] = "kubernetes.io/hostname"
			topologyChanged = true
		}

		if topologyChanged {
			_, err = ki.kubeClient.CoreV1().Services("kube-system").Update(context.TODO(), svc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ki *Initializer) deployOpenYurt() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	converter := &yurtinit.ClusterConverter{
		RootDir:                   dir,
		ComponentsBuilder:         kubeutil.NewBuilder(ki.KubeConfig),
		ClientSet:                 ki.kubeClient,
		CloudNodes:                ki.CloudNodes,
		EdgeNodes:                 ki.EdgeNodes,
		WaitServantJobTimeout:     kubeutil.DefaultWaitServantJobTimeout,
		YurthubHealthCheckTimeout: 2 * time.Minute, // yurtinit.defaultYurthubHealthCheckTimeout
		KubeConfigPath:            ki.KubeConfig,
		YurtManagerImage:          ki.YurtManagerImage,
		NodeServantImage:          ki.NodeServantImage,
		YurthubImage:              ki.YurtHubImage,
		EnableDummyIf:             ki.EnableDummyIf,
	}
	if err := converter.Run(); err != nil {
		klog.Errorf("errors occurred when deploying openyurt components")
		return err
	}
	return nil
}

func (ki *Initializer) loadImagesToKindNodes(images, nodes []string) error {
	for _, image := range images {
		if image == "" {
			// if image == "", it's the responsibility of kind to pull images from registry.
			continue
		}
		if err := ki.operator.KindLoadDockerImage(ki.out, ki.ClusterName, image, nodes); err != nil {
			return err
		}
	}
	return nil
}

func validateKubernetesVersion(ver string) error {
	s := strings.Split(ver, ".")
	var originVer = ver
	if len(s) < 2 || len(s) > 3 {
		return fmt.Errorf("invalid format of kubernetes version: %s", ver)
	}
	if len(s) == 3 {
		// v1.xx.xx
		ver = strings.Join(s[:2], ".")
	}

	if !strings.HasPrefix(ver, "v") {
		ver = fmt.Sprintf("v%s", ver)
	}

	// v1.xx
	if !strutil.IsInStringLst(validKubernetesVersions, ver) {
		return fmt.Errorf("unsupported kubernetes version: %s", originVer)
	}
	return nil
}

func validateOpenYurtVersion(ver string, ignoreError bool) error {
	if !strutil.IsInStringLst(AllValidOpenYurtVersions, ver) && !ignoreError {
		return fmt.Errorf("%s is not a valid openyurt version, all valid versions are %s. If you know what you're doing, you can set --ignore-error",
			ver, strings.Join(AllValidOpenYurtVersions, ","))
	}
	return nil
}

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
