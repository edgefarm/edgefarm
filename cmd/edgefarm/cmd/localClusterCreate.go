// // /*
// // Copyright © 2023 NAME HERE <EMAIL ADDRESS>
// // */
// // package cmd

// // import (
// // 	"fmt"

// // 	"github.com/spf13/cobra"
// // )

// // // locateCreateCmd represents the create command
// // var locateCreateCmd = &cobra.Command{
// // 	Use:   "create",
// // 	Short: "Create a local edgefarm cluster",
// // 	Long:  `Create a local edgefarm cluster`,
// // 	Run: func(cmd *cobra.Command, args []string) {
// // 		fmt.Println("create called")
// // 	},
// // }

// // func init() {
// // 	localCmd.AddCommand(locateCreateCmd)

// // 	// Here you will define your flags and configuration settings.

// // 	// Cobra supports Persistent Flags which will work for this command
// // 	// and all subcommands, e.g.:
// // 	// locateCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

// //		// Cobra supports local flags which will only run when this command
// //		// is called directly, e.g.:
// //		// locateCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// //	}
// package cmd

// import (
// 	goflag "flag"
// 	"os"

// 	"github.com/spf13/cobra"
// 	flag "github.com/spf13/pflag"
// 	"k8s.io/klog/v2"

// 	yurtinit "github.com/edgefarm/edgefarm/pkg/init"
// )

// func NewCreateCommand() *cobra.Command {
// 	cmds := &cobra.Command{
// 		Use:   "create",
// 		Short: "Create a local edgefarm cluster",
// 	}

// 	// add kubeconfig to persistent flags
// 	cmds.PersistentFlags().String("kubeconfig", "", "The path to the kubeconfig file")
// 	cmds.AddCommand(yurtinit.NewInitCMD(os.Stdout))
// 	klog.InitFlags(nil)
// 	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

// 	return cmds
// }

// func init() {
// 	localCmd.AddCommand(NewCreateCommand())
// }

/*
Copyright 2022 The OpenYurt Authors.

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

	yurtinit "github.com/openyurtio/openyurt/test/e2e/cmd/init"

	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
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
		"v0.11.1": {
			"v1.17": "kindest/node:v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00",
			"v1.18": "kindest/node:v1.18.19@sha256:7af1492e19b3192a79f606e43c35fb741e520d195f96399284515f077b3b622c",
			"v1.19": "kindest/node:v1.19.11@sha256:07db187ae84b4b7de440a73886f008cf903fcf5764ba8106a9fd5243d6f32729",
			"v1.20": "kindest/node:v1.20.7@sha256:cbeaf907fc78ac97ce7b625e4bf0de16e3ea725daf6b04f930bd14c67c671ff9",
			"v1.21": "kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6",
		},
		"v0.12.0": {
			"v1.17": "kindest/node:v1.17.17@sha256:e477ee64df5731aa4ef4deabbafc34e8d9a686b49178f726563598344a3898d5",
			"v1.18": "kindest/node:v1.18.20@sha256:e3dca5e16116d11363e31639640042a9b1bd2c90f85717a7fc66be34089a8169",
			"v1.19": "kindest/node:v1.19.16@sha256:81f552397c1e6c1f293f967ecb1344d8857613fb978f963c30e907c32f598467",
			"v1.20": "kindest/node:v1.20.15@sha256:393bb9096c6c4d723bb17bceb0896407d7db581532d11ea2839c80b28e5d8deb",
			"v1.21": "kindest/node:v1.21.10@sha256:84709f09756ba4f863769bdcabe5edafc2ada72d3c8c44d6515fc581b66b029c",
			"v1.22": "siredmar/node:v1.22.7@sha256:d91d1726173bdbc5b5de174884f5afabb602346bb0953b7dc02f27c4d23c1c5e",
			"v1.23": "kindest/node:v1.23.4@sha256:0e34f0d0fd448aa2f2819cfd74e99fe5793a6e4938b328f657c8e3f81ee0dfb9",
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
		Short: "Create a loacl edgefarm cluster",
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
	PortMappings      portMappings
}

type portMappings struct {
	HostApiServerPort    int
	HostNatsPort         int
	HostNatsLeafnodePort int
	HostVaultPort        int
	HostHttpPort         int
	HostHttpsPort        int
}

func newKindOptions() *kindOptions {
	return &kindOptions{
		KindConfigPath:    fmt.Sprintf("%s/kindconfig.yaml", yurtconstantes.TmpDownloadDir),
		WorkerNodesNum:    1,
		EdgeNodesNum:      2,
		ClusterName:       "edgefarm",
		OpenYurtVersion:   "v1.3.0",
		KubernetesVersion: "v1.22",
		UseLocalImages:    false,
		IgnoreError:       true,
		EnableDummyIf:     true,
		DisableDefaultCNI: true,
		CloudNodes:        "edgefarm-control-plane,edgefarm-worker",
		PortMappings: portMappings{
			HostApiServerPort:    6443,
			HostNatsPort:         4222,
			HostNatsLeafnodePort: 7422,
			HostVaultPort:        8200,
			HostHttpPort:         8080,
			HostHttpsPort:        8443,
		},
	}
}

func checkFreePort(port int) bool {
	ln, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", port))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer ln.Close()
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
	if !checkFreePort(o.PortMappings.HostApiServerPort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostApiServerPort)
	}
	if !checkFreePort(o.PortMappings.HostNatsPort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostNatsPort)
	}
	if !checkFreePort(o.PortMappings.HostNatsLeafnodePort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostNatsLeafnodePort)
	}
	if !checkFreePort(o.PortMappings.HostVaultPort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostVaultPort)
	}
	if !checkFreePort(o.PortMappings.HostHttpPort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostHttpPort)
	}
	if !checkFreePort(o.PortMappings.HostHttpsPort) {
		return fmt.Errorf("port %d is already used", o.PortMappings.HostHttpsPort)
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
		PortMappings:      o.PortMappings,
	}
}

func addFlags(flagset *pflag.FlagSet, o *kindOptions) {
	// flagset.StringVar(&o.KindConfigPath, "kind-config-path", o.KindConfigPath,	"Specify the path where the kind config file will be generated.")
	// flagset.IntVar(&o.WorkerNodesNum, "worker-node-num", o.WorkerNodesNum,	"Specify the worker node number of the kind cluster.")
	flagset.IntVar(&o.EdgeNodesNum, "edge-node-num", o.EdgeNodesNum, "Specify the edge node number of the kind cluster.")
	// flagset.StringVar(&o.ClusterName, "cluster-name", o.ClusterName,	"The cluster name of the new-created kind cluster.")
	// flagset.StringVar(&o.CloudNodes, "cloud-nodes", "", 	"Comma separated list of cloud nodes. The control-plane will always be cloud node. If no cloud node specified, the control-plane node will be the only one cloud node.")
	// flagset.StringVar(&o.OpenYurtVersion, "openyurt-version", o.OpenYurtVersion, 	"The version of openyurt components.")
	// flagset.StringVar(&o.KubernetesVersion, "kubernetes-version", o.KubernetesVersion, 	"The version of kubernetes that the openyurt cluster is based on.")
	// flagset.BoolVar(&o.UseLocalImages, "use-local-images", o.UseLocalImages, 	"If set, local images stored by docker will be used first.")
	flagset.StringVar(&o.KubeConfig, "kube-config", o.KubeConfig, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	// flagset.BoolVar(&o.IgnoreError, "ignore-error", o.IgnoreError, 	"Ignore error when using openyurt version that is not officially released.")
	// flagset.BoolVar(&o.EnableDummyIf, "enable-dummy-if", o.EnableDummyIf, 	"Enable dummy interface for yurthub component or not. and recommend to set false on mac env")
	// flagset.BoolVar(&o.DisableDefaultCNI, "disable-default-cni", o.DisableDefaultCNI, 	"Disable the default cni of kind cluster which is kindnet. If this option is set, you should check the ready status of pods by yourself after installing your CNI.")
	flagset.IntVar(&o.PortMappings.HostApiServerPort, "host-api-server-port", o.PortMappings.HostApiServerPort, "Specify the port of host api server.")
	flagset.IntVar(&o.PortMappings.HostVaultPort, "host-vault-port", o.PortMappings.HostVaultPort, "Specify the port of vault to be mapped to.")
	flagset.IntVar(&o.PortMappings.HostNatsPort, "host-nats-port", o.PortMappings.HostNatsPort, "Specify the port of nats to be mapped to.")
	flagset.IntVar(&o.PortMappings.HostNatsLeafnodePort, "host-nats-leafnode-port", o.PortMappings.HostNatsLeafnodePort, "Specify the port of nats leafnode to be mapped to.")
	flagset.IntVar(&o.PortMappings.HostHttpPort, "host-http-port", o.PortMappings.HostHttpPort, "Specify the port of http server to be mapped to.")
	flagset.IntVar(&o.PortMappings.HostHttpsPort, "host-https-port", o.PortMappings.HostHttpsPort, "Specify the port of https server to be mapped to.")
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
	PortMappings      portMappings
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
	klog.Info("Start to install kind")
	if err := ki.operator.KindInstall(); err != nil {
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
		"host_api_server_port": fmt.Sprintf("%d", ki.PortMappings.HostApiServerPort),
	})
	if err != nil {
		return err
	}

	// add additional worker entries into kind config file according to NodesNum
	for num := 0; num < ki.WorkerNodesNum; num++ {
		worker, err := tmplutil.SubsituteTemplate(constants.KindWorkerRole, map[string]string{
			"kind_node_image":         ki.NodeImage,
			"host_nats_port":          fmt.Sprintf("%d", ki.PortMappings.HostNatsPort),
			"host_nats_leafnode_port": fmt.Sprintf("%d", ki.PortMappings.HostNatsLeafnodePort),
			"host_vault_port":         fmt.Sprintf("%d", ki.PortMappings.HostVaultPort),
			"host_http_port":          fmt.Sprintf("%d", ki.PortMappings.HostHttpPort),
			"host_https_port":         fmt.Sprintf("%d", ki.PortMappings.HostHttpsPort),
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
	dp, err := ki.kubeClient.AppsV1().Deployments("kube-system").Get(context.TODO(), "coredns", metav1.GetOptions{})
	if err != nil {
		return err
	}

	if dp != nil {
		nodeList, err := ki.kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		} else if nodeList == nil {
			return fmt.Errorf("failed to list nodes")
		}

		if dp.Spec.Replicas == nil || len(nodeList.Items) != int(*dp.Spec.Replicas) {
			replicas := int32(len(nodeList.Items))
			dp.Spec.Replicas = &replicas
		}

		dp.Spec.Template.Spec.HostNetwork = true

		_, err = ki.kubeClient.AppsV1().Deployments("kube-system").Update(context.TODO(), dp, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

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
