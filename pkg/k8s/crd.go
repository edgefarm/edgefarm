package k8s

import (
	"context"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCRDs() (*v1.CustomResourceDefinitionList, error) {
	clientset, err := apiextensionsclientset.NewForConfig(getConfig())
	if err != nil {
		return nil, err
	}

	crdList, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return crdList, nil
}

func CrdExists(name string) (bool, error) {
	crdList, err := ListCRDs()
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

func CrdEstablished(name string) (bool, error) {
	clientset, err := apiextensionsclientset.NewForConfig(getConfig())
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
