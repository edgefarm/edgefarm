package v1alpha1

import (
	"os"

	"gopkg.in/yaml.v2"
)

func (s *Cluster) Export(path string) error {
	j, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, j, 0644)
	if err != nil {
		return err
	}
	return nil
}
