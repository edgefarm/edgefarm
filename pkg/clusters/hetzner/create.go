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

package hetzner

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/edgefarm/edgefarm/pkg/clusters"
	"github.com/edgefarm/edgefarm/pkg/clusters/capi"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	"github.com/fatih/color"
	"k8s.io/client-go/rest"
)

func CreateCluster(config *rest.Config) error {
	if err := capi.CreateCluster(); err != nil {
		return err
	}

	if err := packages.Install(config, packages.ClusterAPIOperatorHetzner); err != nil {
		return err
	}

	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "caph-system", map[string]string{"cluster.x-k8s.io/provider": "infrastructure-hetzner"}, time.Minute*5); err != nil {
		return err
	}

	_, err := k8s.WaitForCrdEstablished(config, "hcloudmachinetemplates.infrastructure.cluster.x-k8s.io", time.Second*30)
	if err != nil {
		return err
	}

	shared.KubeConfig, err = shared.Expand(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return err
	}

	shared.KubeConfigRestConfig, err = clusters.PrepareKubeClient(shared.KubeConfig)
	if err != nil {
		return err
	}

	state, err := state.GetState(shared.StatePath)
	if err != nil {
		return err
	}
	if state.GetNetbirdSetupKey() == "" {
		return fmt.Errorf("netbird setup key not found")
	}

	context := map[string]string{
		"CLUSTER_NAME":                      shared.ClusterConfig.Spec.Hetzner.Name,
		"KUBERNETES_VERSION":                constants.KubernetesVersion,
		"WORKER_MACHINE_COUNT":              fmt.Sprintf("%d", shared.ClusterConfig.Spec.Hetzner.Workers.Count),
		"HCLOUD_REGION":                     shared.ClusterConfig.Spec.Hetzner.HetznerCloudRegion,
		"CONTROL_PLANE_MACHINE_COUNT":       fmt.Sprintf("%d", shared.ClusterConfig.Spec.Hetzner.ControlPlane.Count),
		"HCLOUD_CONTROL_PLANE_MACHINE_TYPE": shared.ClusterConfig.Spec.Hetzner.ControlPlane.MachineType,
		"HCLOUD_WORKER_MACHINE_TYPE":        shared.ClusterConfig.Spec.Hetzner.Workers.MachineType,
		"HCLOUD_SSH_KEY":                    shared.ClusterConfig.Spec.Hetzner.HetznerCloudSSHKeyName,
		"HCLOUD_TOKEN":                      shared.ClusterConfig.Spec.Hetzner.HCloudToken,
		"NETBIRD_DOMAIN":                    b64.StdEncoding.EncodeToString([]byte("netbird.cloud")),
		"NETBIRD_ADMIN_URL":                 b64.StdEncoding.EncodeToString([]byte("https://app.netbird.io:443")),
		"NETBIRD_MANAGEMENT_URL":            b64.StdEncoding.EncodeToString([]byte("https://api.wiretrustee.com:443")),
		"NETBIRD_SETUP_KEY":                 b64.StdEncoding.EncodeToString([]byte(state.GetNetbirdSetupKey())),
	}
	if err := clusters.RenderAndApply("secret hetzner", hetznerSecret, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err := clusters.RenderAndApply("secret netbird", netbirdSecret, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("secret hetznerSSH", hetznerSSHSecret, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("flannel cloud", flannelCloud, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("flannel edge", flannelEdge, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("coredns", coreDNS, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("kube-proxy default", kubeProxyDefault, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("kube-proxy openyurt", kubeProxyOpenYurt, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("cluster", capiTemplate, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply("helm charts for ccm and csi", hetznerCCMCSI, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	// wait for kubeconfig in secret <cluster-name>-kubeconfig
	exists, err := clusters.WaitForKubeconfig(config, fmt.Sprintf("%s-kubeconfig", shared.ClusterConfig.Spec.Hetzner.Name), time.Minute*5)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("kubeconfig not found")
	}
	s, err := k8s.GetSecret(config, fmt.Sprintf("%s-kubeconfig", shared.ClusterConfig.Spec.Hetzner.Name), "default")
	if err != nil {
		return err
	}
	kubeconfig, err := k8s.SecretValue(s, "value")
	if err != nil {
		return err
	}

	p, err := shared.Expand(shared.ClusterConfig.Spec.Hetzner.KubeConfigPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(p, []byte(kubeconfig), 0644)
	if err != nil {
		return err
	}

	shared.CloudKubeConfig = p
	shared.CloudKubeConfigRestConfig = k8s.GetConfigFromKubeconfig(shared.CloudKubeConfig)

	reachable, err := clusters.WaitForAPIServerReachable(shared.CloudKubeConfigRestConfig, time.Minute*10)
	if err != nil {
		return err
	}
	if !reachable {
		return fmt.Errorf("API server not reachable")
	}

	allNodesPresent, err := clusters.WaitForAllNodes(shared.CloudKubeConfigRestConfig, time.Minute*20, shared.ClusterConfig.Spec.Hetzner.ControlPlane.Count, shared.ClusterConfig.Spec.Hetzner.Workers.Count)
	if err != nil {
		return err
	}
	if !allNodesPresent {
		return fmt.Errorf("not all nodes present")
	}

	return nil
}

func ShowGreeting() {
	green := color.New(color.FgHiGreen)
	yellow := color.New(color.FgHiYellow)
	green.Printf("The hetzner cluster has been created.\nRun ")
	yellow.Printf("  $ local-up deploy --config <hetzner-config.yaml>")
	green.Printf(" to deploy EdgeFarm components and its dependencies.\nHave a look at the arguments using '--help'.")
}
