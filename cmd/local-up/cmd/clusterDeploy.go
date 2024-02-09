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

// import (
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/edgefarm/edgefarm/pkg/args"
// 	"github.com/edgefarm/edgefarm/pkg/k8s"
// 	"github.com/edgefarm/edgefarm/pkg/packages"
// 	"github.com/spf13/cobra"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/klog/v2"
// )

// func NewDeployCommand(out io.Writer) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "deploy",
// 		Short: "Deploy components to the local edgefarm cluster",
// 		RunE: func(cmd *cobra.Command, arguments []string) error {
// 			if err := args.EvaluateKubeConfigPath(); err != nil {
// 				fmt.Printf("Error: %v\n", err)
// 				os.Exit(1)
// 			}
// 			klog.Info("Start to prepare kube client")
// 			kubeconfig, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
// 			if err != nil {
// 				klog.Errorf("Failed to build kubeconfig: %v", err)
// 				os.Exit(1)
// 			}
// 			args.KubeConfigRestConfig = kubeconfig
// 			if err := RunDeploy(); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 		Args: cobra.NoArgs,
// 	}
// 	cmd.SetOut(out)
// 	return cmd
// }

// func init() {
// 	localClusterCmd.AddCommand(NewDeployCommand(os.Stdout))
// 	// localDeleteCmd.PersistentFlags().StringVar(&args.KubeConfig, "kube-config", constants.DefaultKubeConfigPath, "Path where the kubeconfig file of new cluster will be stored. The default is ${HOME}/.kube/config.")
// }

// func RunDeploy() error {
// 	l, err := k8s.ListCRDs()
// 	if err != nil {
// 		return err
// 	}
// 	klog.Infof("CRDs: %v", l)
// 	if err := packages.Install(packages.ClusterBootstrapKyverno); err != nil {
// 		return err
// 	}
// 	// klog.Infof("Prepare edge nodes")
// 	// if err := k8s.PrepareEdgeNodes(); err != nil {
// 	// 	return err
// 	// }
// 	// klog.Infof("Deploy cluster initial packages")
// 	// if err := packages.Install(packages.Init); err != nil {
// 	// 	return err
// 	// }

// 	// klog.Infof("Prepare edge nodes")
// 	// if err := k8s.PrepareEdgeNodes(); err != nil {
// 	// 	return err
// 	// }

// 	// klog.Infof("Deploy cluster base packages")
// 	// if err := packages.Install(packages.Base); err != nil {
// 	// 	return err
// 	// }

// 	// klog.Infof("Deploy cluster dependencies packages")
// 	// if err := packages.Install(packages.ClusterDependencies); err != nil {
// 	// 	return err
// 	// }

// 	// klog.Infof("Deploy edgefarm network packages")
// 	// if err := packages.Install(packages.Network); err != nil {
// 	// 	return err
// 	// }

// 	// // klog.Infof("Deploy edgefarm applications packages")
// 	// if err := packages.Install(packages.Monitor); err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }
