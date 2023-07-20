/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// localClusterCmd represents the local command
var localClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage a local edgefarm cluster",
	Long:  `Manage a local edgefarm cluster`,
}

func init() {
	rootCmd.AddCommand(localClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
