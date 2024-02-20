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

	"github.com/edgefarm/edgefarm/pkg/clusters"
	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"k8s.io/klog/v2"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
)

func CreateCluster() error {
	var err error
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

	shared.KubeConfigRestConfig, err = clusters.PrepareKubeClient(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return err
	}
	if err := packages.Install(shared.KubeConfigRestConfig, packages.CertManager); err != nil {
		return err
	}
	if err := packages.Install(shared.KubeConfigRestConfig, packages.ClusterAPIOperator); err != nil {
		return err
	}
	// if err := packages.Install(packages.CapiToArgoClusterOperator); err != nil {
	// 	return err
	// }

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
