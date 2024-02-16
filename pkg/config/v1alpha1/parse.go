package v1alpha1

import (
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
		shared.Ports.HostApiServerPort = c.Spec.Local.ApiServerPort
		shared.Ports.HostNatsPort = c.Spec.Local.NatsPort
		shared.Ports.HostHttpPort = c.Spec.Local.HttpPort
		shared.Ports.HostHttpsPort = c.Spec.Local.HttpsPort
		shared.EdgeNodesNum = c.Spec.Local.VirtualEdgeNodes
	}
	shared.KubeConfig = c.Spec.General.KubeConfigPath
	return nil
}
