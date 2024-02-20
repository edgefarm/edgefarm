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

package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	strutil "github.com/openyurtio/openyurt/pkg/util/strings"
	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"

	"github.com/edgefarm/edgefarm/pkg/clusters"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/deploy"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/k8s/addons"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/validate"
	"github.com/fatih/color"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"tideland.dev/go/wait"
)

func CreateCluster() error {
	o := newKindOptions()
	if err := validatePorts(); err != nil {
		return err
	}
	initializer, err := newKindInitializer(io.Writer(os.Stdout), o.Config())
	if err != nil {
		return err
	}
	return initializer.Run()
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
	}
}

type kindOptions struct {
	KindConfigPath    string
	WorkerNodesNum    int
	EdgeNodesNum      int
	ClusterName       string
	CloudNodes        string
	OpenYurtVersion   string
	KubernetesVersion string
	KubeConfig        string
	IgnoreError       bool
	NodeImage         string
}

func newKindOptions() *kindOptions {
	return &kindOptions{
		WorkerNodesNum:    1,
		EdgeNodesNum:      shared.EdgeNodesNum,
		ClusterName:       shared.ClusterName,
		OpenYurtVersion:   "v1.4.0",
		KubernetesVersion: constants.KubernetesVersion,
		IgnoreError:       true,
		CloudNodes:        fmt.Sprintf("%s-control-plane,%s-worker", shared.ClusterName, shared.ClusterName),
		NodeImage:         "ghcr.io/edgefarm/edgefarm/kind-node:v1.22.7-systemd",
	}
}

func validatePorts() error {

	if !validate.CheckFreePort(shared.Ports.HostApiServerPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostApiServerPort)
	}
	if !validate.CheckFreePort(shared.Ports.HostNatsPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostNatsPort)
	}
	if !validate.CheckFreePort(shared.Ports.HostHttpPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostHttpPort)
	}
	if !validate.CheckFreePort(shared.Ports.HostHttpsPort) {
		return fmt.Errorf("port %d is already used", shared.Ports.HostHttpsPort)
	}
	return nil
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
	shared.KubeConfigRestConfig, err = clusters.PrepareKubeClient(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return err
	}

	if err := k8s.PrepareEdgeNodes(shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if !shared.Args.Skip.CoreDNS {
		klog.Infof("Deploy cluster coredns packages")
		if err := addons.ReplaceCoreDNS(shared.KubeConfigRestConfig); err != nil {
			return err
		}
	}
	if !shared.Args.Skip.Flannel {
		klog.Infof("Deploy cluster flannel packages")
		if err := packages.Install(shared.KubeConfigRestConfig, packages.Flannel); err != nil {
			return err
		}
		if err := WaitForBootstrapConditions(time.Minute * 5); err != nil {
			return err
		}
	}
	if !shared.Args.Skip.KubeProxy {
		klog.Infof("Deploy cluster kube-proxy packages")
		if err := addons.ReplaceKubeProxy(shared.KubeConfigRestConfig); err != nil {
			return err
		}
	}

	if shared.Args.Deploy {
		if err := deploy.Deploy(shared.KubeConfigRestConfig); err != nil {
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

func (ki *Initializer) prepareKindConfigFile() ([]byte, error) {
	kindConfigContent, err := tmplutil.SubsituteTemplate(kindConfigTemplate, map[string]string{
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
		worker, err := tmplutil.SubsituteTemplate(kindWorkerRoleTemplate, map[string]string{
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
		worker, err := tmplutil.SubsituteTemplate(kindEdgeRole, map[string]string{
			"kind_node_image": ki.NodeImage,
		})
		if err != nil {
			return nil, err
		}
		kindConfigContent = strings.Join([]string{kindConfigContent, worker}, "\n")
	}
	fmt.Println(kindConfigContent)
	return []byte(kindConfigContent), nil
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

func WaitForBootstrapConditions(stepTimeout time.Duration) error {
	ticker := wait.MakeExpiringIntervalTicker(time.Second, stepTimeout)

	// Checks for ready state of all nodes
	nodesCondition := func() (bool, error) {
		nodes, err := k8s.GetAllNodes(shared.KubeConfigRestConfig)
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
