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
	"github.com/edgefarm/edgefarm/pkg/rsa"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func CreateCluster(config *rest.Config) error {
	// check if env variable LOCAL_UP_SKIP_CAPI_BOOTSTRAP exists

	if os.Getenv("LOCAL_UP_SKIP_CAPI_BOOTSTRAP") == "true" {
		klog.Infoln("Skipping creating CAPI cluster")
	} else {
		if err := capi.CreateCluster(); err != nil {
			return err
		}
	}

	if err := packages.Install(config, packages.ClusterAPIOperatorHetzner); err != nil {
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

	privPath, err := shared.Expand(shared.ClusterConfig.Spec.Hetzner.SSHPrivateKeyPath)
	if err != nil {
		return err
	}

	pubPath, err := shared.Expand(shared.ClusterConfig.Spec.Hetzner.SSHPublicKeyPath)
	if err != nil {
		return err
	}

	priv, pub, err := rsa.NewRSA(privPath, pubPath)
	if err != nil {
		return err
	}

	context := map[string]string{
		"CLUSTER_NAME":                      shared.ClusterConfig.Spec.Hetzner.Name,
		"KUBERNETES_VERSION":                constants.KubernetesVersion,
		"WORKER_MACHINE_COUNT":              fmt.Sprintf("%d", shared.ClusterConfig.Spec.Hetzner.WorkerMachineCount),
		"HCLOUD_REGION":                     shared.ClusterConfig.Spec.Hetzner.HetznerCloudRegion,
		"CONTROL_PLANE_MACHINE_COUNT":       fmt.Sprintf("%d", shared.ClusterConfig.Spec.Hetzner.ControlPlaneMachineCount),
		"HCLOUD_CONTROL_PLANE_MACHINE_TYPE": shared.ClusterConfig.Spec.Hetzner.HetznerCloudControlPlaneMachineType,
		"HCLOUD_WORKER_MACHINE_TYPE":        shared.ClusterConfig.Spec.Hetzner.HetznerCloudWorkerMachineType,
		"HCLOUD_SSH_KEY":                    shared.ClusterConfig.Spec.Hetzner.HetznerCloudSSHKey,
		"HCLOUD_TOKEN":                      b64.StdEncoding.EncodeToString([]byte(shared.ClusterConfig.Spec.Hetzner.HCloudToken)),
		"HETZNER_ROBOT_USER":                b64.StdEncoding.EncodeToString([]byte(shared.ClusterConfig.Spec.Hetzner.HetznerRobotUser)),
		"HETZNER_ROBOT_PASSWORD":            b64.StdEncoding.EncodeToString([]byte(shared.ClusterConfig.Spec.Hetzner.HetznerRobotPassword)),
		"HETZNER_SSH_PUBLIC_KEY":            b64.StdEncoding.EncodeToString([]byte(pub)),
		"HETZNER_SSH_PRIVATE_KEY":           b64.StdEncoding.EncodeToString([]byte(priv)),
	}
	if err := clusters.RenderAndApply(hetznerSecret, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply(hetznerSSHSecret, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply(capiTemplate, context, shared.KubeConfigRestConfig); err != nil {
		return err
	}

	if err = clusters.RenderAndApply(hetznerCCMCSI, context, shared.KubeConfigRestConfig); err != nil {
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

	return nil
}
