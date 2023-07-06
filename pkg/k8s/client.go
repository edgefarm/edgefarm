package k8s

import (
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getConfig(kubeconfig *string) *rest.Config {
	configPath := ""
	if kubeconfig == nil {
		if home := homedir.HomeDir(); home != "" {
			configPath = filepath.Join(home, ".kube", "config")
		}
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err.Error())
	}
	return config
}

// GetClientset returns a clientset for the current cluster.
func GetClientset(kubeconfig *string) (*kubernetes.Clientset, error) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

// GetClientset returns a clientset for the current cluster.
func GetDynamicClient(kubeconfig *string) (dynamic.Interface, error) {
	client, err := dynamic.NewForConfig(getConfig(kubeconfig))
	if err != nil {
		panic(err)
	}

	return client, nil
}
