package k8s

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
)

func NodePoolResourceExists() (bool, error) {
	client, err := GetClientset()
	if err != nil {
		return false, err
	}

	groupVersion := schema.GroupVersion{
		Group:   "apps.openyurt.io",
		Version: "v1alpha1",
	}
	apiResourceList, err := client.Discovery().ServerResourcesForGroupVersion(groupVersion.String())
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("failed to discover nodepool resource, %v", err)
		return false, err
	} else if apiResourceList == nil {
		return false, nil
	}

	for i := range apiResourceList.APIResources {
		if apiResourceList.APIResources[i].Name == "nodepools" && apiResourceList.APIResources[i].Kind == "NodePool" {
			return true, nil
		}
	}
	return false, nil
}
