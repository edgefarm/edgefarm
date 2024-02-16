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

package shared

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/spf13/pflag"
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
	ConfigPath           string
	EdgeNodesNum         = 2
	ClusterName          = "edgefarm"
	CloudClusterName     = "edgefarm"
)

type SkipFlags struct {
	Deploy               bool
	KubeProxy            bool
	CoreDNS              bool
	Flannel              bool
	CertManager          bool
	Kyverno              bool
	Crossplane           bool
	VaultOperator        bool
	Metacontroller       bool
	Vault                bool
	EdgeFarmCore         bool
	EdgeFarmApplications bool
	EdgeFarmNetwork      bool
	EdgeFarmMonitor      bool
	Ingress              bool
}

type OnlyFlags struct {
	KubeProxy            bool
	CoreDNS              bool
	Flannel              bool
	CertManager          bool
	Kyverno              bool
	Ingress              bool
	Crossplane           bool
	VaultOperator        bool
	Metacontroller       bool
	Vault                bool
	EdgeFarmCore         bool
	EdgeFarmApplications bool
	EdgeFarmNetwork      bool
	EdgeFarmMonitor      bool
}

type ArgsType struct {
	Deploy bool
	Skip   SkipFlags
	Only   OnlyFlags
}

var (
	Args ArgsType
)

func ConvertOnlyToSkip(only OnlyFlags) SkipFlags {
	skip := SkipFlags{
		Deploy:               true,
		KubeProxy:            true,
		CoreDNS:              true,
		Flannel:              true,
		Ingress:              true,
		CertManager:          true,
		Kyverno:              true,
		Crossplane:           true,
		VaultOperator:        true,
		Metacontroller:       true,
		Vault:                true,
		EdgeFarmCore:         true,
		EdgeFarmApplications: true,
		EdgeFarmNetwork:      true,
		EdgeFarmMonitor:      true,
	}

	if only.KubeProxy {
		skip.KubeProxy = false
	}
	if only.CoreDNS {
		skip.CoreDNS = false
	}
	if only.CertManager {
		skip.CertManager = false
	}
	if only.Flannel {
		skip.Flannel = false
	}
	if only.Kyverno {
		skip.Kyverno = false
	}
	if only.Crossplane {
		skip.Crossplane = false
	}
	if only.VaultOperator {
		skip.VaultOperator = false
	}
	if only.Metacontroller {
		skip.Metacontroller = false
	}
	if only.Vault {
		skip.Vault = false
	}
	if only.EdgeFarmCore {
		skip.EdgeFarmCore = false
	}
	if only.EdgeFarmApplications {
		skip.EdgeFarmApplications = false
	}
	if only.EdgeFarmNetwork {
		skip.EdgeFarmNetwork = false
	}
	if only.EdgeFarmMonitor {
		skip.EdgeFarmMonitor = false
	}

	return skip
}

func EvaluateKubeConfigPath() error {
	kubeConfigPath := KubeConfig
	if strings.HasPrefix(kubeConfigPath, "~") {
		expanded, err := Expand(kubeConfigPath)
		if err != nil {
			return err
		}
		KubeConfig = expanded
		return nil
	}

	if kubeConfigPath == "" {
		if home := os.Getenv("HOME"); home != "" {
			kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
			klog.V(1).Infof("--kube-config is not specified, %s will be used.", kubeConfigPath)
		} else {
			return fmt.Errorf("failed to get ${HOME} env when using default kubeconfig path")
		}
	}

	expanded, err := Expand(kubeConfigPath)
	if err != nil {
		return err
	}
	KubeConfig = expanded
	return nil
}

func Expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func AddSharedFlags(flagset *pflag.FlagSet) {
	flagset.StringVar(&KubeConfig, "kube-config", constants.DefaultKubeConfigPath, fmt.Sprintf("Path where the kubeconfig file of new cluster will be stored. The default is %s", constants.DefaultKubeConfigPath))
	flagset.StringVar(&ConfigPath, "config", ConfigPath, "Path to the edgefarm config file.")

}
