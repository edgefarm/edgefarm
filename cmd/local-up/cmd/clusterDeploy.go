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
	"io"
	"os"

	configv1 "github.com/edgefarm/edgefarm/pkg/config/v1alpha1"
	deploy "github.com/edgefarm/edgefarm/pkg/deploy"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func NewDeployCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy components to the edgefarm cluster",
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
			} else {
				c := configv1.NewConfig(configv1.Local)
				err := configv1.Parse(&c)
				if err != nil {
					return err
				}
			}

			if shared.Args.Only.Crossplane ||
				shared.Args.Only.Kyverno ||
				shared.Args.Only.Metacontroller ||
				shared.Args.Only.Vault ||
				shared.Args.Only.VaultOperator ||
				shared.Args.Only.EdgeFarmApplications ||
				shared.Args.Only.EdgeFarmCore ||
				shared.Args.Only.EdgeFarmMonitor ||
				shared.Args.Only.CertManager ||
				shared.Args.Only.Ingress ||
				shared.Args.Only.EdgeFarmNetwork {
				shared.Args.Skip = shared.ConvertOnlyToSkip(shared.Args.Only)
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
			if err := RunDeploy(configv1.ConfigType(shared.ClusterConfig.Spec.Type), config); err != nil {
				return err
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.SetOut(out)
	shared.AddSharedFlags(cmd.Flags())
	deploy.AddFlagsForDeploy(cmd.Flags())
	return cmd
}

func init() {
	localClusterCmd.AddCommand(NewDeployCommand(os.Stdout))
	// localDeleteCmd.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
}

func RunDeploy(t configv1.ConfigType, config *rest.Config) error {
	// todo: distinguish between local and capi clusters
	if err := deploy.Deploy(t, config); err != nil {
		return err
	}
	return nil
}
