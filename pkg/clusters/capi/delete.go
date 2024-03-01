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
	"os"

	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"k8s.io/klog/v2"
)

func DeleteCluster() error {
	ki, err := kindoperator.NewKindOperator(shared.KubeConfig)
	if err != nil {
		klog.Errorf("Error %v", err)
		os.Exit(1)
	}
	return ki.KindDeleteCluster(shared.ClusterName)
}
