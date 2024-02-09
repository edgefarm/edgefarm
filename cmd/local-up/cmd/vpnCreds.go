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

	"github.com/edgefarm/edgefarm/pkg/state"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var localVPNCredsCmd = &cobra.Command{
	Use:   "creds",
	Short: "Get the VPN creds for connecting to netbird.io.",
	Long: `Get the VPN creds for connecting to netbird.io. 
This command will print the setup-key for the VPN. It only works if the VPN is enabled.`,
	RunE: func(cmd *cobra.Command, arguments []string) error {
		state, err := state.GetState()
		if err != nil {
			klog.Errorf("Failed to get state: %v", err)
			klog.Errorln("VPN status: unknown")
			os.Exit(1)
		}
		if !state.IsFullyConfigured() {
			klog.Infoln("VPN status: disabled. Use 'edgefarm local-up vpn enable to enable VPN.")
			os.Exit(0)
		}
		fmt.Printf("Netbird.io VPN setup-key: %s\n", state.GetNetbirdSetupKey())

		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	localVPNCmd.AddCommand(localVPNCredsCmd)
}
