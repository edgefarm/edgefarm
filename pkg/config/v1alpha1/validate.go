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
	switch {
	case c.Spec.Type == Local.String():
		if err = ValidateLocal(&c.Spec.Local); err != nil {
			return err
		}
	case c.Spec.Type == Hetzner.String():
		if err = ValidateNetbird(&c.Spec.Netbird); err != nil {
			return err
		}
		if err = ValidateHetzner(&c.Spec.Hetzner); err != nil {
			return err
		}
	}

	return nil
}

var (
	hetznerCloudRegions = []string{"fsn1", "nbg1", "hel1", "ash", "hil"}
	hetznerMachiens     = []string{"cx11", "cpx11", "cx21", "cpx21", "cx31", "cpx31", "cx41", "cpx41", "cx51", "cpx51"}
)

func ValidateNetbird(c *api.Netbird) error {
	if c.SetupKey == "" {
		return fmt.Errorf("spec.netbird.setupKey is required. Maybe you should run 'local-up vpn preconfigure --netbird-token <token>'")
	}
	return nil
}

func ValidateHetzner(c *api.Hetzner) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.ControlPlaneMachineCount == 0 {
		return fmt.Errorf("controlPlaneMachineCount is required")
	}
	if c.WorkerMachineCount == 0 {
		return fmt.Errorf("workerMachineCount is required")
	}
	if c.HetznerCloudRegion == "" {
		return fmt.Errorf("hetznerCloudRegion is required")
	}

	// check if region is valid
	regionFound := false
	for _, r := range hetznerCloudRegions {
		if c.HetznerCloudRegion == r {
			regionFound = true
			break
		}
	}
	if !regionFound {
		return fmt.Errorf("invalid hetznerCloudRegion: %s", c.HetznerCloudRegion)
	}

	if c.HetznerCloudControlPlaneMachineType == "" {
		return fmt.Errorf("hetznerCloudControlPlaneMachineType is required")
	}

	// check if machine type is valid
	machineTypeFound := false
	for _, m := range hetznerMachiens {
		if c.HetznerCloudControlPlaneMachineType == m {
			machineTypeFound = true
			break
		}
	}
	if !machineTypeFound {
		return fmt.Errorf("invalid hetznerCloudControlPlaneMachineType: %s", c.HetznerCloudControlPlaneMachineType)
	}

	if c.HetznerCloudWorkerMachineType == "" {
		return fmt.Errorf("hetznerCloudWorkerMachineType is required")
	}

	// check if machine type is valid
	machineTypeFound = false
	for _, m := range hetznerMachiens {
		if c.HetznerCloudWorkerMachineType == m {
			machineTypeFound = true
			break
		}
	}
	if !machineTypeFound {
		return fmt.Errorf("invalid hetznerCloudWorkerMachineType: %s", c.HetznerCloudWorkerMachineType)
	}

	if c.HetznerCloudSSHKey == "" {
		return fmt.Errorf("hetznerCloudSSHKey is required")
	}
	if c.HCloudToken == "" {
		return fmt.Errorf("hcloudToken is required")
	}
	if c.HetznerRobotUser == "" {
		return fmt.Errorf("robotUser is required")
	}
	if c.HetznerRobotPassword == "" {
		return fmt.Errorf("robotPassword is required")
	}
	if c.KubeConfigPath == "" {
		return fmt.Errorf("kubeConfigPath is required")
	}
	if c.HetznerCloudSSHKey == "" {
		return fmt.Errorf("hetznerCloudSSHKey is required")
	}

	return nil
}
