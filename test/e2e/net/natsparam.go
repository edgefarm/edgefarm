package net

import (
	"encoding/base64"
	"fmt"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
)

func getNatsParams(nameSpace string, appName string, compName string, networkName string) (creds string, natsUrl string, err error) {
	secret := fmt.Sprintf("%s.%s", appName, compName)
	jsonPath := fmt.Sprintf("jsonpath='{.data.%s\\.creds}'", networkName)
	out := fw.RunKubectlOrDie(kubeConfig, nameSpace, "get", "secret", secret, "-o", jsonPath)
	if out == "" {
		return "", "", fmt.Errorf("no secret found for network %s", networkName)
	}
	fmt.Printf("creds: %s\n", out)
	credsBytes, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return "", "", err
	}
	creds = string(credsBytes)

	out = fw.RunKubectlOrDie(kubeConfig, nameSpace, "get", "secret", "nets-server-info", "-o", "jsonpath='{.data.NATS_ADDRESS}'")
	if out == "" {
		return "", "", fmt.Errorf("no secret found for nats %s", networkName)
	}
	fmt.Printf("natsUrl: %s\n", out)
	natsUrlBytes, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return "", "", err
	}
	natsUrl = string(natsUrlBytes)
	return creds, natsUrl, nil
}
