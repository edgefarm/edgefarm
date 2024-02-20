/*
Copyright © 2024 EdgeFarm Authors

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
	"time"

	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/netbird"
	"github.com/edgefarm/edgefarm/pkg/shared"
	args "github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/edgefarm/edgefarm/pkg/state"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/andanhm/go-prettytime"
	"github.com/jedib0t/go-pretty/v6/table"
)

type NodePeer struct {
	Node      string
	Connected bool
	Online    bool
	IP        string
	LastSeen  string
	Address   string
}

var localVPNStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the VPN status from netbird.io.",
	Long:  `Get the VPN status from netbird.io.`,
	RunE: func(cmd *cobra.Command, arguments []string) error {
		token := ""
		if args.NetbirdToken != "" {
			token = args.NetbirdToken
		} else {
			tmpstate, err := state.GetState()
			if err != nil {
				klog.Errorf("Failed to get state: %v", err)
				os.Exit(1)
			}
			token = tmpstate.GetNetbirdToken()
		}
		if token == "" {
			klog.Errorln("netbird.io private access not set. Please set the netbird.io private access token.")
			os.Exit(1)
		}

		if err := args.EvaluateKubeConfigPath(); err != nil {
			klog.Errorf("Error: %v\n", err)
			os.Exit(1)
		}
		kubeconfig, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
		if err != nil {
			klog.Errorf("Failed to build kubeconfig: %v", err)
			os.Exit(1)
		}
		args.KubeConfigRestConfig = kubeconfig

		state, err := state.GetState()
		if err != nil {
			klog.Errorf("Failed to get state: %v", err)
			klog.Errorln("VPN status: unknown")
			os.Exit(1)
		}
		if !state.IsFullyConfigured() {
			klog.Infoln("VPN status: disabled. Use 'edgefarm local-up vpn enable to enable VPN.")
			os.Exit(0)
		}

		peers, err := netbird.GetGroupPeers(token)
		if err != nil {
			klog.Errorf("Failed to get relevant peers: %v", err)
			os.Exit(1)
		}

		group, err := netbird.GetGroup(token, state.GetNetbirdGroupID())
		if err != nil {
			klog.Errorf("Failed to get group: %v", err)
			os.Exit(1)
		}

		// todo: distinguish between local and capi clusters
		nodes, err := k8s.GetAllNodes(shared.KubeConfigRestConfig)
		if err != nil {
			klog.Errorf("Failed to get nodes: %v", err)
			os.Exit(1)
		}
		nodePeers := []NodePeer{}
		for _, peer := range peers {
			name, connected := func() (string, bool) {
				for _, currentNode := range nodes {
					if currentNode.Name == peer.Hostname {
						return currentNode.Name, true
					}
				}
				return "", false
			}()

			nodePeer := NodePeer{
				Node:      name,
				Connected: connected,
				Online:    peer.Connected,
				IP:        peer.IP,
				Address:   peer.DNSLabel,
				LastSeen:  peer.LastSeen,
			}
			nodePeers = append(nodePeers, nodePeer)
		}

		problems := []string{}
		tips := []string{}
		printTips := false
		nodesMissingInNetbird := []string{}
		for _, node := range nodes {
			found := false
			for _, np := range nodePeers {
				if np.Node == node.Name {
					found = true
				}
			}
			if !found {
				nodesMissingInNetbird = append(nodesMissingInNetbird, node.Name)
				problems = append(problems, fmt.Sprintf("node '%s' missing in group '%s'", node.Name, group.Name))
				printTips = true
				tips = append(tips, fmt.Sprintf("enable VPN on node '%s'", node.Name))
			}
		}
		// iterate over nodePeers to find duplicates (key Node). If found, add an error message to the comment field
		for i, np := range nodePeers {
			for j, np2 := range nodePeers {
				if i != j && np.Node == np2.Node {
					problems = append(problems, fmt.Sprintf("Duplicated netbird peer for node %s. Something is fishy.", np.Node))
					printTips = true
					tips = append(tips, fmt.Sprintf("- check node '%s' if provisioned using the right setup-key", np.Node))
					tips = append(tips, "check netbird.io for duplicates")
					tips = append(tips, "try to re-enable VPN using 'vpn disable and vpn enable'")
				}
			}
		}

		mytable, err := newTable(table.Row{"Node", "VPN Enabled", "Online", "IP", "DNS Label", "Last seen"})
		if err != nil {
			klog.Errorf("Failed to create table: %v", err)
			os.Exit(1)
		}
		for _, np := range nodePeers {
			t, err := time.Parse("2006-01-02T15:04:05Z", np.LastSeen)
			if err != nil {
				klog.Fatal(err)
			}
			ct := prettytime.Format(t)

			mytable.AppendRows([]table.Row{{
				np.Node, boolSymbol(np.Connected), boolSymbol(np.Online),
				np.IP, np.Address, ct},
			})
		}
		for _, n := range nodesMissingInNetbird {
			mytable.AppendRows([]table.Row{{n, "❌", "❌", "", "", "", ""}})
		}
		mytable.Render()
		// unique problems slice
		uniqueProblems := []string{}
		for _, p := range problems {
			unique := true
			for _, up := range uniqueProblems {
				if p == up {
					unique = false
				}
			}
			if unique {
				uniqueProblems = append(uniqueProblems, p)
			}
		}

		uniqueTips := []string{}
		for _, p := range tips {
			unique := true
			for _, up := range uniqueTips {
				if p == up {
					unique = false
				}
			}
			if unique {
				uniqueTips = append(uniqueTips, p)
			}
		}

		if len(uniqueProblems) > 0 {
			klog.Infoln("Problems:")
			for _, p := range uniqueProblems {
				klog.Infoln("- " + p)
			}
		}
		if printTips {
			klog.Infoln("Tips:")
			for _, t := range uniqueTips {
				klog.Infoln("- " + t)
			}
		}

		return nil
	},
	Args: cobra.NoArgs,
}

func boolSymbol(state bool) string {
	if !state {
		return "❌"
	}
	return "✅"
}

func init() {
	localVPNCmd.AddCommand(localVPNStatusCmd)
	localVPNStatusCmd.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
	localVPNStatusCmd.Flags().StringVar(&args.NetbirdToken, "netbird-token", "", "Specify the netbird.io private access token. This is required to connect physical edge nodes.")
}

func newTable(row table.Row) (table.Writer, error) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendSeparator()
	t.AppendHeader(row)
	t.AppendSeparator()
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	return t, nil
}
