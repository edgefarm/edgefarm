package k8s

import (
	"context"

	"github.com/pytimer/k8sutil/apply"
)

func Apply(manifest string) error {
	dynamicClient, err := GetDynamicClient(nil)
	if err != nil {
		return err
	}

	discoveryClient, err := GetDiscoveryClient(nil)
	if err != nil {
		return err
	}

	applyOptions := apply.NewApplyOptions(dynamicClient, discoveryClient)
	if err := applyOptions.Apply(context.TODO(), []byte(manifest)); err != nil {
		return err
	}
	return nil
}
