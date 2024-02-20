package capi

import (
	"os"

	"github.com/edgefarm/edgefarm/pkg/shared"
)

func GetKubeConfig() (string, error) {
	kubeconfig, err := os.ReadFile(shared.ClusterConfig.Spec.General.KubeConfigPath)
	if err != nil {
		return "", err
	}
	return string(kubeconfig), nil
}
