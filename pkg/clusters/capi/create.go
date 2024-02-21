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

package capi

import (
	"fmt"
	"os"
	"time"

	"github.com/edgefarm/edgefarm/pkg/clusters"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/fatih/color"
	"k8s.io/klog/v2"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
)

func CreateCluster() error {
	var err error
	if os.Getenv("LOCAL_UP_SKIP_CAPI_BOOTSTRAP") == "true" {
		klog.Infoln("Skipping creating CAPI cluster")
	} else {
		ki, err := kindoperator.NewKindOperator(shared.ClusterConfig.Spec.General.KubeConfigPath)
		if err != nil {
			klog.Errorf("Error %v", err)
			os.Exit(1)
		}

		config, err := prepareKindConfigFile(shared.ClusterName)
		if err != nil {
			return err
		}

		err = ki.KindCreateClusterWithConfig(config)
		if err != nil {
			return err
		}
	}
	shared.KubeConfigRestConfig, err = clusters.PrepareKubeClient(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return err
	}
	if err := packages.Install(shared.KubeConfigRestConfig, packages.CertManager); err != nil {
		return err
	}

	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "cert-manager", map[string]string{"app.kubernetes.io/instance": "cert-manager"}, time.Minute*5); err != nil {
		return err
	}

	if err := packages.Install(shared.KubeConfigRestConfig, packages.ClusterAPIOperator); err != nil {
		return err
	}

	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "capi-kubeadm-bootstrap-system", map[string]string{"cluster.x-k8s.io/provider": "bootstrap-kubeadm"}, time.Minute*5); err != nil {
		return err
	}
	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "capi-kubeadm-control-plane-system", map[string]string{"cluster.x-k8s.io/provider": "control-plane-kubeadm"}, time.Minute*5); err != nil {
		return err
	}
	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "capi-operator-system", map[string]string{"cluster.x-k8s.io/provider": "capi-operator"}, time.Minute*5); err != nil {
		return err
	}
	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "capi-system", map[string]string{"cluster.x-k8s.io/provider": "cluster-api"}, time.Minute*5); err != nil {
		return err
	}
	if err := k8s.WaitForDeploymentOrError(shared.KubeConfigRestConfig, "caaph-system", map[string]string{"cluster.x-k8s.io/provider": "helm"}, time.Minute*5); err != nil {
		return err
	}

	return nil
}

func prepareKindConfigFile(name string) ([]byte, error) {
	kindConfigContent, err := tmplutil.SubsituteTemplate(kindConfigTemplate, map[string]string{
		"cluster_name": name,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(kindConfigContent)
	return []byte(kindConfigContent), nil
}

func ShowGreeting() {
	green := color.New(color.FgHiGreen)
	green.Printf("The cluster-api-bootstrap cluster has been created. Proceeding with the cloud cluster creation\n")
}
