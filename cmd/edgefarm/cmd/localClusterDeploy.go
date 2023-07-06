/*
Copyright 2022 The OpenYurt Authors.

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
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/edgefarm/edgefarm/pkg/packages"
)

func NewDeployCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy components to the local edgefarm cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Run(); err != nil {
				return err
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.SetOut(out)
	return cmd
}

func init() {
	localClusterCmd.AddCommand(NewDeployCommand(os.Stdout))
}

func Run() error {
	// klog.Infof("Deploy cluster initial packages")
	// if err := packages.Install(packages.Init); err != nil {
	// 	return err
	// }

	// klog.Infof("Prepare edge nodes")
	// if err := k8s.PrepareEdgeNodes(); err != nil {
	// 	return err
	// }

	// klog.Infof("Deploy cluster base packages")
	// if err := packages.Install(packages.Base); err != nil {
	// 	return err
	// }

	// klog.Infof("Deploy cluster dependencies packages")
	// if err := packages.Install(packages.Dependencies); err != nil {
	// 	return err
	// }

	klog.Infof("Deploy edgefarm network packages")
	if err := packages.Install(packages.Network); err != nil {
		return err
	}

	return nil
}
