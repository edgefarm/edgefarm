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
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func GetSecret(kubeconfig *rest.Config, name, namespace string) (*v1.Secret, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func SecretValue(secret *v1.Secret, key string) (string, error) {
	val, exists := secret.Data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in secret", key)
	}
	return string(val), nil
}
