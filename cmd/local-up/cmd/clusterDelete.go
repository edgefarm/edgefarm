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
	"io"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/edgefarm/edgefarm/pkg/kindoperator"
	"github.com/edgefarm/edgefarm/pkg/netbird"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
)

var (
	override       bool
	cleanupNetbird bool
)

// localDeleteCmd represents the localDelete command
func NewClusterDeleteCommand(out io.Writer) *cobra.Command {
	localDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the local edgefarm cluster",
		Run: func(cmd *cobra.Command, arguments []string) {
			klog.Infoln("Deleting local edgefarm cluster")
			if err := args.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			klog.Info("Start to prepare kube client")
			kubeconfig, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
			if err != nil {
				klog.Errorf("Failed to build kubeconfig: %v", err)
				os.Exit(1)
			}
			args.KubeConfigRestConfig = kubeconfig

			ki, err := kindoperator.NewKindOperator(args.KubeConfig)
			if err != nil {
				klog.Errorf("Error %v", err)
				os.Exit(1)
			}

			doit := false
			if override {
				doit = true
			} else {
				var err error
				input := confirmation.New("Are you sure to delete the local edgefarm cluster?", confirmation.Yes)
				doit, err = input.RunPrompt()
				if err != nil {
					klog.Errorf("Error %v", err)
					os.Exit(1)
				}
			}
			if doit {
				state, err := state.GetState()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				if state.IsFullyConfigured() && cleanupNetbird {
					klog.Infoln("netbird.io: cleanup")
					err := netbird.DisableVPN(false, true, true, true, true)
					if err != nil {
						klog.Errorf("Error %v", err)
						os.Exit(1)
					}
				}
				err = ki.KindDeleteCluster("edgefarm")
				if err != nil {
					klog.Errorf("Error %v", err)
					os.Exit(1)
				}
			} else {
				klog.Infoln("Aborted")
			}
		},
	}
	localDeleteCmd.SetOut(out)
	shared.AddSharedFlags(localDeleteCmd.Flags())
	addDeleteArgs(localDeleteCmd.Flags())
	return localDeleteCmd
}

func init() {
	localClusterCmd.AddCommand(NewClusterDeleteCommand(os.Stdout))
	// localDeleteCmd.Flags().BoolVarP(&override, "yes", "y", false, "Override confirmation prompt")
	// localDeleteCmd.Flags().BoolVarP(&cleanupNetbird, "cleanup-netbird", "c", true, "Cleanup netbird.io resources")
}

func addDeleteArgs(flagset *pflag.FlagSet) {
	flagset.BoolVarP(&override, "yes", "y", false, "Override confirmation prompt")
	flagset.BoolVarP(&cleanupNetbird, "cleanup-netbird", "c", true, "Cleanup netbird.io resources")

}
