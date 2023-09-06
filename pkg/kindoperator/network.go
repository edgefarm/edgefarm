package kindoperator

import (
	"context"
	"net"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func getSubnetAndGateway(subnet string) (string, string, error) {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", "", err
	}

	s := ipNet.String()
	tokens := strings.Split(s, ".")
	gw := strings.Join(tokens[:len(tokens)-1], ".") + ".1"

	return s, gw, nil
}

func networkExists(name string) (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, err
	}

	_, err = cli.NetworkInspect(context.Background(), name, types.NetworkInspectOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func deleteNetwork(name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	err = cli.NetworkRemove(context.Background(), name)
	if err != nil {
		return err
	}
	return nil
}

func createNetwork(name string, subnet string) error {
	s, gw, err := getSubnetAndGateway(subnet)
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	_, err = cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Scope:          "local",
		EnableIPv6:     false,
		IPAM: &network.IPAM{
			Driver:  "",
			Options: map[string]string{},
			Config: []network.IPAMConfig{
				{
					Subnet:  s,
					Gateway: gw,
				},
			},
		},
		Internal:   false,
		Attachable: false,
		Ingress:    false,
		ConfigOnly: false,
		Options: map[string]string{
			"com.docker.network.bridge.enable_ip_masquerade": "true",
			"com.docker.network.driver.mtu":                  "1500",
		},
		Labels: map[string]string{},
	})
	if err != nil {
		return err
	}
	return nil
}
