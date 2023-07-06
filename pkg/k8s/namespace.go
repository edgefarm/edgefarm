package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespace(name string) (*v1.Namespace, error) {
	clientset, err := GetClientset(nil)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
}

func CreateNamespace(name string) (*v1.Namespace, error) {
	clientset, err := GetClientset(nil)
	if err != nil {
		return nil, err
	}

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
}

func CreateNamespaceIfNotExist(name string) (*v1.Namespace, error) {
	create := false
	ns, err := GetNamespace(name)
	if err != nil {
		if errors.IsNotFound(err) {
			create = true
		} else {
			return nil, err
		}
	}
	if create {
		ns, err = CreateNamespace(name)
	}
	return ns, err
}
