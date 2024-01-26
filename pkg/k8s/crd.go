package k8s

import (
	"context"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCRDs() (*v1beta1.CustomResourceDefinitionList, error) {
	clientset, err := apiextensionsclientset.NewForConfig(getConfig())
	if err != nil {
		return nil, err
	}

	return clientset.ApiextensionsV1beta1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
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
