package k8s

import (
	"context"

	"github.com/pytimer/k8sutil/apply"
	"k8s.io/client-go/rest"
)

func Apply(config *rest.Config, manifest string) error {
	dynamicClient, err := GetDynamicClient(config)
	if err != nil {
		return err
	}

	discoveryClient, err := GetDiscoveryClient(config)
	if err != nil {
		return err
	}

	applyOptions := apply.NewApplyOptions(dynamicClient, discoveryClient)
	if err := applyOptions.Apply(context.TODO(), []byte(manifest)); err != nil {
		return err
	}
	return nil
}
