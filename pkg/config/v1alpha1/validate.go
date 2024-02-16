package v1alpha1

import (
	"fmt"

	api "github.com/edgefarm/edgefarm/apis/config/v1alpha1"
)

func ValidateType(arg string) error {
	if arg == "" {
		return fmt.Errorf("cluster type is required")
	}
	clusterTypeFound := false
	for _, t := range ValidClusterTypes {
		if arg == t.String() {
			clusterTypeFound = true
			break
		}
	}
	if !clusterTypeFound {
		return fmt.Errorf("invalid cluster type: %s", arg)
	}
	return nil
}

func ValidateGeneral(c *api.General) error {
	if c.KubeConfigPath == "" {
		return fmt.Errorf("kubeConfigPath is required")
	}
	return nil
}

func ValidateLocal(c *api.Local) error {
	if c.ApiServerPort == 0 {
		return fmt.Errorf("apiServerPort is required")
	}
	if c.NatsPort == 0 {
		return fmt.Errorf("natsPort is required")
	}
	if c.HttpPort == 0 {
		return fmt.Errorf("httpPort is required")
	}
	if c.HttpsPort == 0 {
		return fmt.Errorf("httpsPort is required")
	}
	return nil
}

func Validate(c *api.Cluster) error {
	if c.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}
	if c.Kind == "" {
		return fmt.Errorf("kind is required")
	}
	if c.APIVersion != "config.edgefarm.io/v1alpha1" {
		return fmt.Errorf("invalid apiVersion: %s", c.APIVersion)
	}
	if c.Kind != "Cluster" {
		return fmt.Errorf("invalid kind: %s", c.Kind)
	}

	err := ValidateType(c.Spec.Type)
	if err != nil {
		return err
	}
	if err = ValidateGeneral(&c.Spec.General); err != nil {
		return err
	}
	if c.Spec.Type == Local.String() {
		if err = ValidateLocal(&c.Spec.Local); err != nil {
			return err
		}
	}
	return nil
}
