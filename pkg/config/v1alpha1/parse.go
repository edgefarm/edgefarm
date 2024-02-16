package v1alpha1

import (
	"fmt"
	"os"

	api "github.com/edgefarm/edgefarm/apis/config/v1alpha1"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"gopkg.in/yaml.v2"
)

func Load(path string) (*api.Cluster, error) {
	path, err := shared.Expand(path)
	if err != nil {
		return nil, err
	}

	str, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cluster := api.Cluster{}
	err = yaml.Unmarshal(str, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, err
}

func Parse(c *api.Cluster) error {
	if c.Spec.Type == Local.String() {
		shared.ClusterName = c.Spec.Local.Name
		shared.Ports.HostApiServerPort = c.Spec.Local.ApiServerPort
		shared.Ports.HostNatsPort = c.Spec.Local.NatsPort
		shared.Ports.HostHttpPort = c.Spec.Local.HttpPort
		shared.Ports.HostHttpsPort = c.Spec.Local.HttpsPort
		shared.EdgeNodesNum = c.Spec.Local.VirtualEdgeNodes
	} else if c.Spec.Type == Hetzner.String() {
		shared.ClusterName = fmt.Sprintf("%s-bootstrap", c.Spec.Hetzner.Name)
		shared.CloudClusterName = c.Spec.Hetzner.Name
	}
	shared.KubeConfig = c.Spec.General.KubeConfigPath
	return nil
}
