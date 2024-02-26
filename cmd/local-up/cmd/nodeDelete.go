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

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	configv1 "github.com/edgefarm/edgefarm/pkg/config/v1alpha1"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	nodeNameDeleteNode string
)

func validateDeleteNode(config *rest.Config) error {
	if nodeNameDeleteNode == "" {
		return errors.New("name must be specified")
	}
	exists, err := k8s.NodeExists(config, nodeNameDeleteNode)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	err = k8s.ValidatePhysicalNodeName(nodeNameDeleteNode)
	if err != nil {
		return err
	}
	return nil
}

func NewNodeDeleteCommand(config *rest.Config, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a edge node from the cluster.",
		Long:  "Deletes a edge node from the cluster. It also gives instructions on how to unprovision the the device.",
		RunE: func(cmd *cobra.Command, arguments []string) error {
			if shared.ConfigPath != "" {
				c, err := configv1.Load(shared.ConfigPath)
				if err != nil {
					return err
				}
				err = configv1.Parse(c)
				if err != nil {
					return err
				}
			}

			switch shared.ClusterType {
			case configv1.Local.String():
				shared.KubeConfig = shared.ClusterConfig.Spec.General.KubeConfigPath
			case configv1.Hetzner.String():
				shared.KubeConfig = shared.ClusterConfig.Spec.Hetzner.KubeConfigPath
			}

			if err := shared.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			klog.Info("Start to prepare kube client")
			config, err := clientcmd.BuildConfigFromFlags("", shared.KubeConfig)
			if err != nil {
				klog.Errorf("Failed to build kubeconfig: %v", err)
				os.Exit(1)
			}

			shared.KubeConfigRestConfig = config
			if err := validateDeleteNode(config); err != nil {
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
	shared.AddSharedFlags(cmd.Flags())
	return cmd
}

func init() {
	nodeDeleteCommand := NewNodeDeleteCommand(shared.KubeConfigRestConfig, os.Stdout)
	nodeCmd.AddCommand(nodeDeleteCommand)
	nodeDeleteCommand.Flags().StringVarP(&nodeNameDeleteNode, "name", "n", "", "The name of the node to delete. Must be one of the self-provisioned nodes.")

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
	klog.Infof("Delete nodepool for node %s", nodeNameDeleteNode)
	err := k8s.DeleteNodepool(shared.KubeConfigRestConfig, nodeNameDeleteNode)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.Infof("Nodepool %s not found", nodeNameDeleteNode)
			goto proceed
		}
		return err
	}
proceed:
	klog.Infof("Delete node %s", nodeNameDeleteNode)
	err = k8s.DeleteNode(shared.KubeConfigRestConfig, nodeNameDeleteNode)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.Infof("Node %s not found", nodeNameDeleteNode)
			goto proceed2
		}
		return err
	}
proceed2:
	instructionsDeleteNode()

	return nil
}
