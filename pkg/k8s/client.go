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

package k8s

import (
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfigFromKubeconfig(kubeconfig string) *rest.Config {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func getConfig(config *rest.Config) *rest.Config {
	if config != nil {
		return config
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

// GetClientset returns a clientset for the current cluster.
func GetClientset(config *rest.Config) (*kubernetes.Clientset, error) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(getConfig(config))
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetClientset returns a clientset for the current cluster.
func GetDynamicClient(config *rest.Config) (dynamic.Interface, error) {
	client, err := dynamic.NewForConfig(getConfig(config))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDiscoveryClient(config *rest.Config) (*discovery.DiscoveryClient, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(getConfig(config))
	if err != nil {
		return nil, err
	}
	return discoveryClient, nil
}
