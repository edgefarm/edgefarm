package v1alpha1

import (
	api "github.com/edgefarm/edgefarm/apis/config/v1alpha1"
)

func NewConfig(t ConfigType) api.Cluster {
	c := api.Cluster{}
	api.SetDefaultsCluster(&c)
	api.SetDefaultsGeneral(&c.Spec.General)
	switch {
	case t == Local:
		c.Spec.Type = Local.String()
		api.SetDefaultsLocal(&c.Spec.Local)
	case t == Hetzner:
		c.Spec.Type = Hetzner.String()
		api.SetDefaultsHetzner(&c.Spec.Hetzner)
		api.SetDefaultNetbird(&c.Spec.Netbird)
	}

	return c
}
