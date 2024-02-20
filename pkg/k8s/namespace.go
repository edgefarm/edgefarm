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
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func GetNamespace(kubeconfig *rest.Config, name string) (*v1.Namespace, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
}

func CreateNamespace(kubeconfig *rest.Config, name string) (*v1.Namespace, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
}

func CreateNamespaceIfNotExist(kubeconfig *rest.Config, name string) (*v1.Namespace, error) {
	create := false
	ns, err := GetNamespace(kubeconfig, name)
	if err != nil {
		if errors.IsNotFound(err) {
			create = true
		} else {
			return nil, err
		}
	}
	if create {
		ns, err = CreateNamespace(kubeconfig, name)
	}
	return ns, err
}
