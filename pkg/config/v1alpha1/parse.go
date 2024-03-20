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
	shared.ClusterConfig = c
	shared.ClusterType = c.Spec.Type
	if c.Spec.Type == Local.String() {
		shared.ClusterName = c.Metadata.Name
		shared.Ports.HostApiServerPort = c.Spec.Local.ApiServerPort
		shared.Ports.HostNatsPort = c.Spec.Local.NatsPort
		shared.Ports.HostHttpPort = c.Spec.Local.HttpPort
		shared.Ports.HostHttpsPort = c.Spec.Local.HttpsPort
		shared.EdgeNodesNum = c.Spec.Local.VirtualEdgeNodes
	} else if c.Spec.Type == Hetzner.String() {
		shared.ClusterName = fmt.Sprintf("%s-bootstrap", c.Metadata.Name)
		shared.CloudClusterName = c.Metadata.Name
	}
	shared.KubeConfig = c.Spec.General.KubeConfigPath
	shared.StatePath = c.Spec.General.StatePath
	return nil
}
