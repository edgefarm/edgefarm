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

	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/packages"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplaceKubeProxy deletes the kube-proxy deployment and replaces it with a Helm chart
func ReplaceKubeProxy() error {
	clientset, err := k8s.GetClientset(nil)
	if err != nil {
		return err
	}

	// // Check if kube-proxy daemonset exists, if true, delete it
	// ds, err := clientset.AppsV1().DaemonSets("kube-system").Get(context.Background(), "kube-proxy", metav1.GetOptions{})
	// if err != nil {
	// 	if apierrors.IsNotFound(err) {
	// 		ds = nil
	// 	} else {
	// 		return err
	// 	}
	// 	return err
	// }

	// if ds != nil {
	err = clientset.AppsV1().DaemonSets("kube-system").Delete(context.Background(), "kube-proxy", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	// }
	err = clientset.CoreV1().ConfigMaps("kube-system").Delete(context.Background(), "kube-proxy", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = clientset.CoreV1().ServiceAccounts("kube-system").Delete(context.Background(), "kube-proxy", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	if err := packages.Install(packages.ClusterBootstrapKubeProxy); err != nil {
		return err
	}

	// ticker := wait.MakeExpiringIntervalTicker(time.Second, time.Second*60)

	// condition := func() (bool, error) {
	// 	pods, err := k8s.GetPods("kube-system", "k8s-app=kube-proxy-default")
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	for _, pod := range pods {
	// 		if pod.Status.Phase != v1.PodRunning {
	// 			return false, nil
	// 		}
	// 	}
	// 	return true, nil
	// }
	// err = wait.Poll(context.Background(), ticker, condition)
	// if err != nil {
	// 	return fmt.Errorf("install kube-proxy-default: %s", err.Error())
	// }

	// ticker = wait.MakeExpiringIntervalTicker(time.Second, time.Second*60)
	// condition = func() (bool, error) {
	// 	pods, err := k8s.GetPods("kube-system", "k8s-app=kube-proxy-openyurt")
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	for _, pod := range pods {
	// 		if pod.Status.Phase != v1.PodRunning {
	// 			return false, nil
	// 		}
	// 	}
	// 	return true, nil
	// }
	// err = wait.Poll(context.Background(), ticker, condition)
	// if err != nil {
	// 	return fmt.Errorf("install kube-proxy-openyurt: %s", err.Error())
	// }

	return nil
}
