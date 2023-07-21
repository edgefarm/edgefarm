/*
Copyright © 2023 EdgeFarm authors

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

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print version information

This command displays version information for the local-up tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("local-up version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
