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

package addons

import (
	"context"
	"time"

	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/packages"
	"tideland.dev/go/wait"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplaceCoreDNS deletes the CoreDNS deployment and replace it with a DaemonSet
func ReplaceCoreDNS() error {
	clientset, err := k8s.GetClientset()
	if err != nil {
		return err
	}

	err = clientset.AppsV1().Deployments("kube-system").Delete(context.Background(), "coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.CoreV1().ConfigMaps("kube-system").Delete(context.Background(), "coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.CoreV1().Services("kube-system").Delete(context.Background(), "kube-dns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.CoreV1().ServiceAccounts("kube-system").Delete(context.Background(), "coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.RbacV1().ClusterRoleBindings().Delete(context.Background(), "system:coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.RbacV1().ClusterRoles().Delete(context.Background(), "system:coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	if err := packages.Install(packages.CoreDNS); err != nil {
		return err
	}

	ticker := wait.MakeExpiringIntervalTicker(time.Second, time.Second*60)

	condition := func() (bool, error) {
		pods, err := k8s.GetPods("kube-system", "k8s-app=kube-dns")
		if err != nil {
			return false, err
		}
		for _, pod := range pods {
			if pod.Status.Phase != v1.PodRunning {
				return false, nil
			}
		}
		return true, nil
	}
	err = wait.Poll(context.Background(), ticker, condition)
	if err != nil {
		return err
	}

	return nil
}
