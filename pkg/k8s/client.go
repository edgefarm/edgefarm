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
	"path/filepath"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getConfig(kubeconfig *string) *rest.Config {
	configPath := ""
	if kubeconfig == nil {
		if home := homedir.HomeDir(); home != "" {
			configPath = filepath.Join(home, ".kube", "config")
		}
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err.Error())
	}
	return config
}

// GetClientset returns a clientset for the current cluster.
func GetClientset(kubeconfig *string) (*kubernetes.Clientset, error) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetClientset returns a clientset for the current cluster.
func GetDynamicClient(kubeconfig *string) (dynamic.Interface, error) {
	client, err := dynamic.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDiscoveryClient(kubeconfig *string) (*discovery.DiscoveryClient, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(getConfig(kubeconfig))
	if err != nil {
		return nil, err
	}
	return discoveryClient, nil
}
