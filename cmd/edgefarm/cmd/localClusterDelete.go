/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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
	kubeConfig string = "${HOME}/.kube/config"
	override   bool
)

// localDeleteCmd represents the localDelete command
var localDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a local edgefarm cluster",
	Long:  `Delete a local edgefarm cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting local edgefarm cluster")
		ki := kindoperator.NewKindOperator("", kubeConfig)
		if err := ki.KindInstall(); err != nil {
			fmt.Println(err)
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
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
		if doit {
			ki.KindDeleteCluster(os.Stdout, "edgefarm")
		} else {
			fmt.Println("Aborted")
		}
	},
}

func init() {
	localClusterCmd.AddCommand(localDeleteCmd)
	localDeleteCmd.PersistentFlags().StringVar(&kubeConfig, "kube-config", kubeConfig, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	// Here you will define your flags and configuration settings.
	localDeleteCmd.Flags().BoolVarP(&override, "yes", "y", false, "Override confirmation prompt")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
