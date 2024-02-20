package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/s0rg/retry"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func ListCRDs(kubeconfig *rest.Config) (*v1.CustomResourceDefinitionList, error) {
	clientset, err := apiextensionsclientset.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		return nil, err
	}

	crdList, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return crdList, nil
}

func CrdExists(kubeconfig *rest.Config, name string) (bool, error) {
	crdList, err := ListCRDs(kubeconfig)
	if err != nil {
		return false, err
	}

	for _, crd := range crdList.Items {
		if crd.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func CrdEstablished(kubeconfig *rest.Config, name string) (bool, error) {
	clientset, err := apiextensionsclientset.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		return false, err
	}
	crd, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	for _, cond := range crd.Status.Conditions {
		if cond.Type == v1.Established && cond.Status == v1.ConditionTrue {
			return true, nil
		}
	}
	return false, nil
}

func WaitForCrdEstablished(kubeconfig *rest.Config, name string, timeout time.Duration) (bool, error) {
	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single(fmt.Sprintf("Waiting for CRD %s to be established", name),
		func() error {
			est, err := CrdEstablished(kubeconfig, name)
			if err != nil {
				return err
			}
			if est {
				return nil
			}
			return nil
		}); err != nil {
		return false, err
	}
	return true, nil
}
