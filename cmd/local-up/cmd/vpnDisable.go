// /*
// Copyright © 2023 EdgeFarm Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package cmd

import (
	"os"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/netbird"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var localVPNDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables netbird.io VPN capabilities for the local edgefarm cluster",
	Long: `Disables netbird.io VPN capabilities for the local edgefarm cluster.
You won't be able to run physical edge nodes using the local edgefarm cluster.`,
	RunE: func(cmd *cobra.Command, arguments []string) error {
		if err := args.EvaluateKubeConfigPath(); err != nil {
			klog.Errorf("Error: %v\n", err)
			os.Exit(1)
		}

		if args.NetbirdToken == "" {
			klog.Infof("netbird.io private access not set. Using cached token.\n")
		}

		if err := netbird.DisableVPN(true, true, true, true); err != nil {
			return err
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	localVPNCmd.AddCommand(localVPNDisableCmd)
	localVPNDisableCmd.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	localVPNDisableCmd.Flags().StringVar(&args.NetbirdToken, "netbird-token", "", "Specify the netbird.io private access token. This is required to connect physical edge nodes.")
}
