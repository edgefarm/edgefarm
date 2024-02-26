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
	"os"

	"github.com/edgefarm/edgefarm/pkg/netbird"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	unconfigureVPN bool
)

func NewVPNPreconfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preconfigure",
		Short: "Preconfigures VPN capabilities using netbird.io for a cloud based edgefarm cluster",
		Long: `Preconfigures VPN capabilities using netbird.io for a cloud based edgefarm cluster.
This enables you to join physical edge nodes to the cloud based edgefarm cluster.`,
		RunE: func(cmd *cobra.Command, arguments []string) error {
			if args.NetbirdToken == "" {
				klog.Errorln("Error: netbird.io private access token must be specified")
				os.Exit(1)
			}
			if unconfigureVPN {
				err := netbird.UnPreconfigure()
				if err != nil {
					return err
				}
				klog.Infoln("VPN unconfiguration completed successfully")
				return nil
			} else {
				key, err := netbird.Preconfigure()
				if err != nil {
					return err
				}
				klog.Infoln("VPN preconfiguration completed successfully")
				klog.Infof("Netbird.io setup-key: %s\n", key)
				return nil
			}
		},
		Args: cobra.NoArgs,
	}
	vpnPreconfigureFlags(cmd.Flags())
	return cmd
}

func init() {
	localVPNCmd.AddCommand(NewVPNPreconfigureCommand())
}

func vpnPreconfigureFlags(flagset *pflag.FlagSet) {
	flagset.StringVar(&args.NetbirdToken, "netbird-token", "", "Specify the netbird.io private access token. This is required to connect physical edge nodes.")
	flagset.BoolVar(&unconfigureVPN, "delete", false, "Deletes pre-configured VPN from netbird.io.")
}
