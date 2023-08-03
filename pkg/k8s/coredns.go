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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PatchCoreDNSDeployment patches the CoreDNS deployment to contain the edgefarm.io/NoSchedule toleration
func PatchCoreDNSDeployment() error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	deployment, err := clientset.AppsV1().Deployments("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	if err != nil {
		return err
	}

	tolerations := deployment.Spec.Template.Spec.Tolerations
	if tolerations == nil {
		tolerations = []v1.Toleration{}
	}

	tolerations = append(tolerations, v1.Toleration{
		Key:    "edgefarm.io",
		Effect: v1.TaintEffectNoSchedule,
	})

	deployment.Spec.Template.Spec.Tolerations = tolerations
	_, err = clientset.AppsV1().Deployments("kube-system").Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
