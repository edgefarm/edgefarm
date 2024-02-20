package k8s

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func PollForConfigMap(kubeconfig *rest.Config, namespace, name string, timeout time.Duration) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for configMap %s/%s", namespace, name)
		default:
			_, err = clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if !errors.IsNotFound(err) {
					return nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func GetConfigMapValue(kubeconfig *rest.Config, namespace, name, key string) (string, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return "", err
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("configMap %s/%s not found", namespace, name)
		}
		return "", err
	}

	value := cm.Data[key]
	if value == "" {
		return "", fmt.Errorf("key not found in ConfigMap")
	}

	return value, nil
}

func UpdateConfigMapValue(kubeconfig *rest.Config, namespace, name, key, value string) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return fmt.Errorf("configMap not found")
	}

	cm.Data[key] = value
	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
