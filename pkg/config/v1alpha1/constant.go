package v1alpha1

const (
	Local   ConfigType = "local"
	Hetzner ConfigType = "hetzner"
)

var (
	ValidClusterTypes = []ConfigType{Local, Hetzner}
)

type ConfigType string

func (c ConfigType) String() string {
	return string(c)
}
