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
	"io"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kubectllogs "k8s.io/kubectl/pkg/cmd/logs"
)

func GetPods(namespace string, selector string) ([]v1.Pod, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, err
	}

	listOptions := metav1.ListOptions{
		LabelSelector: selector,
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func PrintPodLog(client kubeclientset.Interface, pod *corev1.Pod, w io.Writer) error {
	klog.Infof("start to print logs for pod(%s/%s):", pod.Namespace, pod.Name)
	req := client.CoreV1().Pods(pod.GetNamespace()).GetLogs(pod.Name, &corev1.PodLogOptions{})
	if err := kubectllogs.DefaultConsumeRequest(req, w); err != nil {
		klog.Errorf("failed to print logs for pod(%s/%s), %v", pod.Namespace, pod.Name, err)
		return err
	}

	return nil
}
