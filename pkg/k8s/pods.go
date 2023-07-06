package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(namespace string, selector string) ([]v1.Pod, error) {
	clientset, err := GetClientset(nil)
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
