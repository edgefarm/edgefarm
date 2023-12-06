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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	nodeNameDeleteNode string
)

func validateDeleteNode() error {
	if nodeNameDeleteNode == "" {
		return errors.New("name must be specified")
	}
	exists, err := k8s.NodeExists(nodeNameDeleteNode)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("node does not exist")
	}
	err = k8s.ValidatePhysicalNodeName(nodeNameDeleteNode)
	if err != nil {
		return err
	}
	return nil
}

func NewNodeDeleteCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a edge node from the cluster.",
		Long:  "Deletes a edge node from the cluster. It also gives instructions on how to unprovision the the device.",
		RunE: func(cmd *cobra.Command, arguments []string) error {
			if err := args.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if err := validateDeleteNode(); err != nil {
				return err
			}

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
	nodeDeleteCommand := NewNodeDeleteCommand(os.Stdout)
	nodeCmd.AddCommand(nodeDeleteCommand)
	nodeDeleteCommand.Flags().StringVarP(&nodeNameDeleteNode, "name", "n", "", "The name of the node to delete. Must be one of the self-provisioned nodes.")
	nodeDeleteCommand.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")

}

func instructionsDeleteNode() {
	green := color.New(color.FgHiGreen)
	yellow := color.New(color.FgHiYellow)

	green.Printf("To unprovision a physical edge node follow the instructions.\n")
	green.Printf("If joined using yurtadm run on the edge node:\n")
	yellow.Println("yurtadm reset -f")
	yellow.Println("rm -rf /etc/cni/net.d")
	yellow.Println("ip link set cni0 down")
	yellow.Println("ip link delete cni0")
	yellow.Println("ip link set flannel.1 down")
	yellow.Println("ip link delete flannel.1")
	yellow.Println("ip link set yurthub-dummy0 down")
	yellow.Println("ip link delete yurthub-dummy0")
	yellow.Println("systemctl restart docker")
}

func Run() error {
	klog.Infof("Delete node %s", nodeNameDeleteNode)
	err := k8s.DeleteNode(nodeNameDeleteNode)
	if err != nil {
		return err
	}

	klog.Infof("Delete nodepool for node %s", nodeNameDeleteNode)
	err = k8s.DeleteNodepool(nodeNameDeleteNode)
	if err != nil {
		return err
	}

	instructionsDeleteNode()

	return nil
}
