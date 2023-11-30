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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteServiceAccount(name, namespace string) error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	return clientset.CoreV1().ServiceAccounts(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func DeleteClusterRole(name string) error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	return clientset.RbacV1().ClusterRoles().Delete(context.Background(), name, metav1.DeleteOptions{})
}

func DeleteClusterRoleBinding(name string) error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	return clientset.RbacV1().ClusterRoleBindings().Delete(context.Background(), name, metav1.DeleteOptions{})
}
