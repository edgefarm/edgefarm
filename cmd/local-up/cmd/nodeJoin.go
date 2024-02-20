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
	"time"

	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/k8s/tokens"
	"github.com/edgefarm/edgefarm/pkg/netbird"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	"github.com/fatih/color"
	"github.com/hako/durafmt"
	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

var (
	nodeNameJoinNode string
	TTL              string
	defaultTTL       string = "24h"
)

func validateJoinNode(config *rest.Config) error {
	state, err := state.GetState()
	if err != nil {
		return err
	}
	if state.GetNetbirdSetupKey() == "" {
		return errors.New("cluster is not VPN enabled. Please run 'local-up vpn enable' first")
	}

	if nodeNameJoinNode == "" {
		return errors.New("name must be specified")
	}
	exists, err := k8s.NodeExists(config, nodeNameJoinNode)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("node already exists")
	}

	return nil
}

func NewNodeJoinCommand(config *rest.Config, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Join a new node to the cluster.",
		Long:  "Join a new node to the cluster by creating a new kubernetes node and giving instructions on how to join it to the cluster.",
		RunE: func(cmd *cobra.Command, arguments []string) error {
			if err := args.EvaluateKubeConfigPath(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if err := validateJoinNode(config); err != nil {
				return err
			}

			if err := RunJoinNode(); err != nil {
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
	nodeJoinCommand := NewNodeJoinCommand(shared.KubeConfigRestConfig, os.Stdout)
	nodeCmd.AddCommand(nodeJoinCommand)
	nodeJoinCommand.Flags().StringVarP(&nodeNameJoinNode, "name", "n", "", "A unique name of the node to join. Must not be the same as an existing node.")
	nodeJoinCommand.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	nodeJoinCommand.PersistentFlags().StringVar(&TTL, "ttl", defaultTTL, "Define the TTL of the bootstrap token.")
}

func instructionsJoinNode(token string, ttl string) error {
	state, err := state.GetState()
	if err != nil {
		return err
	}
	routingPeer, err := netbird.RoutingPeerIP(state.GetNetbirdToken())
	if err != nil {
		return err
	}

	green := color.New(color.FgHiGreen)
	greenBold := color.New(color.FgHiGreen, color.Bold)
	yellow := color.New(color.FgHiYellow)

	green.Printf("Here is some information you need to join a physical edge node to this cluster.\n\n")
	greenBold.Println("VPN:")
	green.Println("Unless you already connected the physical node to netbird.io VPN, you need to connect it to the VPN first.")
	green.Println("")
	green.Printf("Use can use this setup-key ")
	yellow.Printf("%s", state.GetNetbirdSetupKey())
	green.Printf(" to connect to netbird.io VPN.\n\n")
	greenBold.Println("Kubernetes:")
	green.Printf("Ensure that the ")
	yellow.Printf("/etc/hosts")
	green.Printf(" file on your physical edge node contains the following entry:\n")
	yellow.Printf("%s edgefarm-control-plane\n", routingPeer)
	green.Println("")
	green.Printf("Use this token ")
	yellow.Printf("%s", token)
	green.Printf(" to join the cluster. You have ")
	yellow.Printf("%s", ttl)
	green.Println(" to join the cluster before this token expires.")
	yellow.Println("")
	green.Println("If you experience any problems, please consult the documentation at ")
	green.Println("https://edgefarm.github.io/edgefarm/ or file an issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md")
	return nil
}

func RunJoinNode() error {
	klog.Infof("Adding empty node resource for %s", nodeNameJoinNode)

	klog.Infof("Adding nodepool for node %s", nodeNameJoinNode)
	nodepoolManifest, err := tmplutil.SubsituteTemplate(constants.NodepoolManifest, map[string]string{
		"name": nodeNameJoinNode,
	})
	if err != nil {
		return err
	}

	client, err := k8s.GetClientset(shared.KubeConfigRestConfig)
	if err != nil {
		return err
	}
	ttlDuration, err := time.ParseDuration(TTL)
	if err != nil {
		return err
	}

	token, err := tokens.GenerateBootstrapToken(client, ttlDuration)
	if err != nil {
		return err
	}

	err = k8s.Apply(shared.KubeConfigRestConfig, nodepoolManifest)
	if err != nil {
		return err
	}

	duraTTL := durafmt.Parse(ttlDuration)
	instructionsJoinNode(token, duraTTL.String())

	return nil
}
