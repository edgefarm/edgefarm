package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/s0rg/retry"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func DeleteCluster(name, namespace string, config *rest.Config) error {
	dynamicClient, err := GetDynamicClient(config)
	if err != nil {
		return err
	}

	clusterRes := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "clusters",
	}

	clusterObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cluster.x-k8s.io/v1beta1",
			"kind":       "Cluster",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
		},
	}

	err = dynamicClient.Resource(clusterRes).Namespace(namespace).Delete(context.TODO(), clusterObject.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetCluster(name, namespace string, config *rest.Config) (*unstructured.Unstructured, error) {
	dynamicClient, err := GetDynamicClient(config)
	if err != nil {
		return nil, err
	}

	clusterRes := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "clusters",
	}

	return dynamicClient.Resource(clusterRes).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func WaitForClusterDeleted(name, namespace string, timeout time.Duration, config *rest.Config) (bool, error) {
	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single(fmt.Sprintf("Waiting for Cluster deleted %s/%s", namespace, name),
		func() error {
			c, err := GetCluster(name, "default", config)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			if c != nil {
				return fmt.Errorf("cluster exists")
			}
			return nil
		}); err != nil {
		return false, err
	}
	return true, nil
}
