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
	"fmt"
	"time"

	"github.com/edgefarm/edgefarm/pkg/clusters"
	"github.com/edgefarm/edgefarm/pkg/clusters/capi"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"k8s.io/apimachinery/pkg/api/errors"
)

func DeleteCluster() error {
	var err error
	shared.KubeConfigRestConfig, err = clusters.PrepareKubeClient(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return err
	}
	if err = k8s.DeleteCluster(shared.CloudClusterName, "default", shared.KubeConfigRestConfig); err != nil {
		if errors.IsNotFound(err) {
			goto proceed
		}
		return err
	}
	{
		deleted, err := k8s.WaitForClusterDeleted(shared.CloudClusterName, "default", time.Minute*5, shared.KubeConfigRestConfig)
		if err != nil {
			return err
		}
		if !deleted {
			return fmt.Errorf("cluster %s not deleted", shared.CloudClusterName)
		}
	}
proceed:

	if err := capi.DeleteCluster(); err != nil {
		return err
	}

	return nil
}
