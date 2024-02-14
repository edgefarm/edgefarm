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

package kindoperator

import (
	"os"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"

	constants "github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/shared"
)

const (
	KindClusterName   = "edgefarm"
	KindNetworkName   = "edgefarm"
	KindNetworkSubnet = "172.254.0.0/16"
)

type KindOperator struct {
	kubeconfigPath string
	logger         log.Logger
}

func NewKindOperator(kubeconfigPath string) (*KindOperator, error) {
	path := constants.DefaultKubeConfigPath
	if kubeconfigPath != "" {
		path = kubeconfigPath
	}
	path, err := shared.Expand(path)
	if err != nil {
		return nil, err
	}

	return &KindOperator{
		kubeconfigPath: path,
		logger:         NewLogger(os.Stdout, 0),
	}, nil
}

func (k *KindOperator) KindCreateClusterWithConfig(config []byte) error {
	exists, err := networkExists(KindNetworkName)
	if err != nil {
		return err
	}
	if !exists {
		err := createNetwork(KindNetworkName, KindNetworkSubnet)
		if err != nil {
			return err
		}
	}

	provider := cluster.NewProvider(cluster.ProviderWithLogger(k.logger))
	err = os.Setenv("KIND_EXPERIMENTAL_DOCKER_NETWORK", KindNetworkName)
	if err != nil {
		return err
	}
	options := []cluster.CreateOption{

		cluster.CreateWithRawConfig(config),
		cluster.CreateWithRetain(true),
		cluster.CreateWithWaitForReady(0),
		cluster.CreateWithKubeconfigPath(k.kubeconfigPath),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
	}

	err = provider.Create(KindClusterName, options...)
	if err != nil {
		return err
	}

	return nil
}

func (k *KindOperator) KindDeleteCluster(name string) error {
	provider := cluster.NewProvider()
	err := provider.Delete(KindClusterName, k.kubeconfigPath)
	if err != nil {
		return err
	}
	exists, err := networkExists(KindNetworkName)
	if err != nil {
		return err
	}
	if exists {
		err = deleteNetwork(KindNetworkName)
		if err != nil {
			return err
		}
	}
	return nil
}
