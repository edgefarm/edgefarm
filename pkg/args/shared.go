/*
Copyright © 2023 EdgeFarm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package args

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type PortMappings struct {
	HostApiServerPort int
	HostNatsPort      int
	HostHttpPort      int
	HostHttpsPort     int
}

var (
	Ports = PortMappings{
		HostApiServerPort: 6443,
		HostNatsPort:      4222,
		HostHttpPort:      80,
		HostHttpsPort:     443,
	}
	NetbirdToken    string
	NetbirdSetupKey string

	KubeConfig           string
	KubeConfigRestConfig *rest.Config
)

func EvaluateKubeConfigPath() error {
	kubeConfigPath := KubeConfig
	if kubeConfigPath == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
			klog.V(1).Infof("--kube-config is not specified, %s will be used.", kubeConfigPath)
		} else {
			return fmt.Errorf("failed to get ${HOME} env when using default kubeconfig path")
		}
	}

	expanded, err := expand(kubeConfigPath)
	if err != nil {
		return err
	}
	KubeConfig = expanded
	return nil
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
