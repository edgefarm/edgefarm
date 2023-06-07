/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// localDeleteCmd represents the localDelete command
var localNodeJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Joins a physical node to the local edgefarm cluster",
	Long:  `Joins a physical node to the local edgefarm cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: create node.yaml, nodepool.yaml, print instructions to download yurtadm and print join command")
	},
}

func init() {
	localNodeCmd.AddCommand(localNodeJoinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
