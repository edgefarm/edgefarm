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
	"io"
	"os"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/route"
	"github.com/fatih/color"
	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	nodeNameJoinNode string
)

func validateJoinNode() error {
	if nodeNameJoinNode == "" {
		return errors.New("name must be specified")
	}
	exists, err := k8s.NodeExists(nodeNameJoinNode)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("node already exists")
	}

	return nil
}

func NewNodeJoinCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Join a new node to the cluster.",
		Long:  "Join a new node to the cluster by creating a new kubernetes node and giving instructions on how to join it to the cluster.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateJoinNode(); err != nil {
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
	nodeJoinCommand := NewNodeJoinCommand(os.Stdout)
	nodeCmd.AddCommand(nodeJoinCommand)
	nodeJoinCommand.Flags().StringVarP(&nodeNameJoinNode, "name", "n", "", "A unique name of the node to join. Must not be the same as an existing node.")
}

func instructionsJoinNode() {
	green := color.New(color.FgHiGreen)
	yellow := color.New(color.FgHiYellow)

	r, err := route.GetRoute(args.Interface)
	if err != nil {
		klog.Fatalf("Failed to get route: %v", err)
		panic(err)
	}

	green.Println("To join a physical edge node to this cluster, you need to ensure a few things:")
	green.Printf("1. Ensure that the ")
	yellow.Printf("/etc/hosts")
	green.Printf(" file contain the following entry:\n")
	yellow.Printf("   %s edgefarm-control-plane\n", r.IP)
	green.Println("")
	green.Println("2. Ensure that the cgroup2 is enabled for your system. To check run:")
	yellow.Println("   $ grep cgroup /proc/filesystems")
	green.Println("   If cgroup2 is enabled, you should see a line like this:")
	yellow.Println("   nodev   cgroup2")
	green.Println("   If cgroup2 is not enabled, enable it. Follow the instructions for your")
	green.Println("   distribution and architecture.")
	green.Println("")
	green.Println("3. Ensure that the cgroup root kubelet and kubelet.slice exist. To create run:")
	yellow.Println("   $ mkdir -p /sys/fs/cgroup/kubelet")
	yellow.Println("   $ mkdir -p /sys/fs/cgroup/kubelet.slice")
	// yellow.Println(`   $ sed -e 's/ / +/g' -e 's/^/+/' <"/sys/fs/cgroup/kubelet/cgroup.controllers" >"/sys/fs/cgroup/kubelet/cgroup.subtree_control"`)
	// yellow.Println(`   $ sed -e 's/ / +/g' -e 's/^/+/' <"/sys/fs/cgroup/kubelet.slice/cgroup.controllers" >"/sys/fs/cgroup/kubelet.slice/cgroup.subtree_control"`)
	green.Println("")
	green.Println("4. Ensure that docker is installed, running and uses the cgroupfs driver.")
	green.Println("   To check run:")
	yellow.Println("   $ docker info | grep -i cgroupfs")
	green.Println("   You should see a line like this:")
	yellow.Println("   Cgroup Driver: cgroupfs")
	green.Println("   If you don't see this line, you need to configure docker to use the cgroupfs")
	green.Println("    driver. https://docs.docker.com/engine/reference/commandline/dockerd/#configure-cgroup-driver")
	green.Println("   For example, you can configure docker to use the cgroupfs driver by creating ")
	green.Println("   a file at /etc/docker/daemon.json with the following contents:   ")
	yellow.Println("   $ cat /etc/docker/daemon.json ")
	yellow.Println("   {")
	yellow.Println(`     "exec-opts": ["native.cgroupdriver=cgroupfs"]`)
	yellow.Println("   }")
	green.Println("")
	green.Println("5. Ensure that you have yurtadm installed on your the edge node. To install visit:")
	yellow.Printf("   https://github.com/openyurtio/openyurt/releases/tag/v1.2.2")
	green.Printf(" and download the\n   yurtadm binary for your platform.\n")
	green.Println("")
	green.Println("Raspberry PI note: Please make sure that these values")
	yellow.Printf("   `cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory`")
	green.Printf(" are present\n   in your cmdline.txt file on your boot partition.\n")
	green.Println("")
	green.Println("Once you have ensured that the above requirements are met, you can")
	green.Println("   join this node to the cluster by running the following command on the node:")
	yellow.Println(`   $ yurtadm join edgefarm-control-plane:6443  --token="abcdef.0123456789abcdef" --node-type=edge --discovery-token-unsafe-skip-ca-verification --v=5`)
	yellow.Println("")
	green.Println("If you experience any problems, please consult the documentation at ")
	green.Println("https://edgefarm.github.io/edgefarm/ or file a issue at https://github.com/edgefarm/edgefarm/issues/new?template=question.md")
}

func RunJoinNode() error {
	klog.Infof("Adding empty node resource for %s", nodeNameJoinNode)
	nodeManifest, err := tmplutil.SubsituteTemplate(constants.NodeManifest, map[string]string{
		"name": nodeNameJoinNode,
	})
	if err != nil {
		return err
	}

	klog.Infof("Adding nodepool for node %s", nodeNameJoinNode)
	nodepoolManifest, err := tmplutil.SubsituteTemplate(constants.NodepoolManifest, map[string]string{
		"name": nodeNameJoinNode,
	})
	if err != nil {
		return err
	}

	err = k8s.Apply(nodeManifest)
	if err != nil {
		return err
	}

	err = k8s.Apply(nodepoolManifest)
	if err != nil {
		return err
	}

	instructionsJoinNode()

	return nil
}
