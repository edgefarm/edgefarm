package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/s0rg/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func DeploymentsReadyByLabels(kubeconfig *rest.Config, namespace string, labels map[string]string) (bool, error) {
	client, err := GetClientset(kubeconfig)
	if err != nil {
		return false, err
	}
	depList, err := client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(metav1.SetAsLabelSelector(labels))})
	if err != nil {
		return false, err
	}
	for _, dep := range depList.Items {
		if dep.Status.ReadyReplicas != dep.Status.Replicas {
			return false, nil
		}
	}
	return true, nil
}

func WaitForDeploymentReady(kubeconfig *rest.Config, namespace string, labels map[string]string, timeout time.Duration) (bool, error) {
	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single(fmt.Sprintf("Waiting for deployments in namespace %s to be ready", namespace),
		func() error {
			ready, err := DeploymentsReadyByLabels(kubeconfig, namespace, labels)
			if err != nil {
				return err
			}
			if !ready {
				return fmt.Errorf("deployments in ns %s not ready", namespace)
			}
			return nil
		}); err != nil {
		return false, err
	}
	return true, nil
}

func WaitForDeploymentOrError(config *rest.Config, namespace string, labels map[string]string, timeout time.Duration) error {
	ready, err := WaitForDeploymentReady(config, namespace, labels, timeout)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("deployment in ns %s not ready", namespace)
	}
	return nil
}
