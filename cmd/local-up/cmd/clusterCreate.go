/*
Copyright © 2024 EdgeFarm Authors

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

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/edgefarm/edgefarm/pkg/clusters/hetzner"
	"github.com/edgefarm/edgefarm/pkg/clusters/local"
	configv1 "github.com/edgefarm/edgefarm/pkg/config/v1alpha1"
	deploy "github.com/edgefarm/edgefarm/pkg/deploy"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	stringsx "github.com/icza/gox/stringsx"
)

var (
	generateConfig bool
)

func NewCreateCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the local edgefarm cluster",
		RunE: func(cmd *cobra.Command, arguments []string) error {
			if generateConfig {
				err := configv1.ValidateType(shared.ClusterType)
				if err != nil {
					return err
				}
				if shared.ClusterType != "" {
					c := configv1.NewConfig(configv1.ConfigType(shared.ClusterType))
					str, err := yaml.Marshal(c)
					if err != nil {
						return err
					}
					fmt.Println(string(str))
					os.Exit(0)
				}
			}
			if shared.ConfigPath != "" {
				c, err := configv1.Load(shared.ConfigPath)
				if err != nil {
					return err
				}
				err = configv1.Parse(c)
				if err != nil {
					return err
				}
			}

			_, err := state.GetState()
			if err != nil {
				return err
			}
			if err := shared.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if shared.Args.Only.CoreDNS || shared.Args.Only.Flannel || shared.Args.Only.KubeProxy {
				shared.Args.Skip = shared.ConvertOnlyToSkip(shared.Args.Only)
			}
			err = validatepProcFS()
			if err != nil {
				return err
			}
			switch {
			case shared.ClusterType == configv1.Local.String():
				err := local.CreateCluster()
				if err != nil {
					return err
				}
				local.ShowGreeting()

			case shared.ClusterType == configv1.Hetzner.String():
				err := hetzner.CreateCluster(shared.KubeConfigRestConfig)
				if err != nil {
					return err
				}
				hetzner.ShowGreeting()
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.SetOut(out)
	shared.AddSharedFlags(cmd.Flags())
	addFlagsForCreate(cmd.Flags())
	deploy.AddFlagsForDeploy(cmd.Flags())
	return cmd
}

func init() {
	localClusterCmd.AddCommand(NewCreateCommand(os.Stdout))
}

func validatepProcFS() error {
	if runtime.GOOS == "linux" {
		notGood := false
		f, err := os.ReadFile("/proc/sys/fs/inotify/max_user_instances")
		if err != nil {
			return err
		}

		v, err := strconv.Atoi(stringsx.Clean(string(f)))
		if err != nil {
			return err
		}
		if v < 512 {
			notGood = true
			klog.Errorln("the value of /proc/sys/fs/inotify/max_user_instances must be greater or equal than 512")
		}

		f, err = os.ReadFile("/proc/sys/fs/inotify/max_user_watches")
		if err != nil {
			return err
		}

		v, err = strconv.Atoi(stringsx.Clean(string(f)))
		if err != nil {
			return err
		}

		if v < 524288 {
			notGood = true
			klog.Errorln("the value of /proc/sys/fs/inotify/max_user_watches must be greater or equal than 524288")
		}

		if notGood {
			return errors.New("follow https://kind.sigs.k8s.io/docs/user/known-issues/#pod-errors-due-to-too-many-open-files to fix this issue")
		}
	}

	return nil
}

func addFlagsForCreate(flagset *pflag.FlagSet) {
	flagset.IntVar(&shared.EdgeNodesNum, "edge-node-num", shared.EdgeNodesNum, "Specify the edge node number of the kind cluster.")
	flagset.BoolVar(&generateConfig, "generate-config", false, "Generates a config file and exit.")
	flagset.StringVar(&shared.ClusterType, "type", "local", fmt.Sprintf("Config type to generate. Valid values are %s", func() string {
		res := ""
		for _, t := range configv1.ValidClusterTypes {
			res += t.String() + ","
		}
		return res[:len(res)-1]
	}()))
	flagset.IntVar(&shared.Ports.HostApiServerPort, "host-api-server-port", shared.Ports.HostApiServerPort, "Specify the port of host api server.")
	flagset.IntVar(&shared.Ports.HostNatsPort, "host-nats-port", shared.Ports.HostNatsPort, "Specify the port of nats to be mapped to.")
	flagset.IntVar(&shared.Ports.HostHttpPort, "host-http-port", shared.Ports.HostHttpPort, "Specify the port of http server to be mapped to.")
	flagset.IntVar(&shared.Ports.HostHttpsPort, "host-https-port", shared.Ports.HostHttpsPort, "Specify the port of https server to be mapped to.")
	flagset.BoolVar(&shared.Args.Deploy, "deploy", false, "Deploy the cluster after creation.")
	flagset.BoolVar(&shared.Args.Skip.CoreDNS, "skip-coredns", false, "Skip installing CoreDNS. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Skip.KubeProxy, "skip-kube-proxy", false, "Skip installing kube-proxy. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Skip.Flannel, "skip-flannel", false, "Skip installing flannel. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.CoreDNS, "only-coredns", false, "Only install CoreDNS. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.Flannel, "only-flannel", false, "Only install flannel. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
	flagset.BoolVar(&shared.Args.Only.KubeProxy, "only-kube-proxy", false, "Only install kube-proxy. WARNING: HERE BE DRAGONS. Make sure your kube context is correct! Use at your own risk.")
}
