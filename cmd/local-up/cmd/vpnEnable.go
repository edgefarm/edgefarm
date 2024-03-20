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
	"fmt"
	"os"

	apiv1 "github.com/edgefarm/edgefarm/apis/config/v1alpha1"
	configv1 "github.com/edgefarm/edgefarm/pkg/config/v1alpha1"
	"github.com/edgefarm/edgefarm/pkg/netbird"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func NewLocalVPNEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enabled VPN capabilities using netbird.io for the local edgefarm cluster",
		Long: `Enabled VPN capabilities using netbird.io for the local edgefarm cluster.
This enables you to join physical edge nodes to the local edgefarm cluster.`,
		RunE: func(cmd *cobra.Command, arguments []string) error {
			var config *apiv1.Cluster
			if shared.ConfigPath != "" {
				c, err := configv1.Load(shared.ConfigPath)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				err = configv1.Parse(c)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				config = c
			} else {
				c := configv1.NewConfig(configv1.Local)
				err := configv1.Parse(&c)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				config = &c
			}
			if config.Spec.Type != configv1.Local.String() {
				return fmt.Errorf("only clusters of type 'local' is supported for this command")
			}

			if err := args.EvaluateKubeConfigPath(); err != nil {
				klog.Errorf("Error: %v\n", err)
				os.Exit(1)
			}
			klog.Info("Start to prepare kube client")
			kubeconfig, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
			if err != nil {
				klog.Errorf("Failed to build kubeconfig: %v", err)
				os.Exit(1)
			}
			args.KubeConfigRestConfig = kubeconfig

			if args.NetbirdToken == "" {
				klog.Errorln("Error: netbird.io private access token must be specified")
				os.Exit(1)
			}

			if err := netbird.EnableVPN(config.Metadata.Name); err != nil {
				return err
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	shared.AddSharedFlags(cmd.Flags())
	addVPNEnableArgs(cmd.Flags())
	return cmd
}

func init() {
	localVPNCmd.AddCommand(NewLocalVPNEnableCmd())
}

func addVPNEnableArgs(flagset *pflag.FlagSet) {
	flagset.StringVar(&args.NetbirdToken, "netbird-token", "", "Specify the netbird.io private access token. This is required to connect physical edge nodes.")
}
