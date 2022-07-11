package net

import (
	"encoding/base64"
	"fmt"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/util/wait"
)

// startupNatsSubscriber waits until nats credentials are available and
// returns a new NatsSubscriber
func startupNatsSubscriber(nameSpace string, appName string, compName string, networkName string) (*NatsSubscriber, error) {
	var creds string
	var natsUrl string
	err := wait.PollImmediate(time.Second*5, dsPollTimeout, func() (bool, error) {
		var err error
		creds, natsUrl, err = getNatsParams(nameSpace, appName, compName, networkName)
		if err != nil {
			g.GinkgoWriter.Printf("Error getting nats params: %v\n", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	g.GinkgoWriter.Printf("Nats params: %s %s\n", creds, natsUrl)

	sub, err := NewNatsSubscriber(natsUrl, creds, "EXPORT.>", "e2e-consumer",
		fmt.Sprintf("%s_%s", net1Name, streamName))
	return sub, err
}

func getNatsParams(nameSpace string, appName string, compName string, networkName string) (creds string, natsUrl string, err error) {
	secret := fmt.Sprintf("%s.%s", appName, compName)
	jsonPath := fmt.Sprintf("jsonpath={.data.%s\\.creds}", networkName)
	out, err := fw.RunKubectl(kubeConfig, nameSpace, "get", "secret", secret, "-o", jsonPath)
	if err != nil {
		return "", "", err
	}
	if out == "" {
		return "", "", fmt.Errorf("no secret found for network %s", networkName)
	}
	credsBytes, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return "", "", err
	}
	creds = string(credsBytes)

	out, err = fw.RunKubectl(kubeConfig, nameSpace, "get", "secret", "nats-server-info", "-o", "jsonpath={.data.NATS_ADDRESS}")
	if err != nil {
		return "", "", err
	}
	if out == "" {
		return "", "", fmt.Errorf("no secret found for nats %s", networkName)
	}
	natsUrlBytes, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return "", "", err
	}
	natsUrl = string(natsUrlBytes)
	return creds, natsUrl, nil
}
