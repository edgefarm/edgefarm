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

package cmd

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/spf13/cobra"

	"github.com/edgefarm/edgefarm/pkg/kindoperator"
)

var (
	kubeConfig string = "~/.kube/config"
	override   bool
)

// localDeleteCmd represents the localDelete command
var localDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the local edgefarm cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting local edgefarm cluster")
		ki := kindoperator.NewKindOperator(kubeConfig)
		doit := false
		if override {
			doit = true
		} else {
			var err error
			input := confirmation.New("Are you sure to delete the local edgefarm cluster?", confirmation.Yes)
			doit, err = input.RunPrompt()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
		if doit {
			err := ki.KindDeleteCluster("edgefarm")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Aborted")
		}
	},
}

func init() {
	localClusterCmd.AddCommand(localDeleteCmd)
	localDeleteCmd.PersistentFlags().StringVar(&kubeConfig, "kube-config", kubeConfig, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	localDeleteCmd.Flags().BoolVarP(&override, "yes", "y", false, "Override confirmation prompt")
}
