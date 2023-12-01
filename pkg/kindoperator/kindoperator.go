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

package kindoperator

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"

	strutil "github.com/openyurtio/openyurt/pkg/util/strings"
)

var (
	// enableGO111MODULE will be used when 1.13 <= go version <= 1.16
	defaultKubeConfigPath = "${HOME}/.kube/config"

	validKindVersions = []string{
		"v0.12.0",
	}
)

type KindOperator struct {
	// kubeconfig represents where to store the kubeconfig of the new created cluster.
	kindCMDPath    string
	kubeconfigPath string

	// KindOperator will use this function to get command.
	// For the convenience of stub in unit tests.
	execCommand func(string, ...string) *exec.Cmd
}

func NewKindOperator(kindCMDPath string, kubeconfigPath string) *KindOperator {
	path := defaultKubeConfigPath
	if kubeconfigPath != "" {
		path = kubeconfigPath
	}
	return &KindOperator{
		kubeconfigPath: path,
		kindCMDPath:    kindCMDPath,
		execCommand:    exec.Command,
	}
}

// GetKindPath returns the path of kind command.
func (k *KindOperator) GetKindPath() error {
	if k.kindCMDPath != "" {
		return nil
	}

	kindPath, err := findKindPath()
	if err != nil {
		klog.Infof("no kind tool is found, so try to install. %v", err)
	} else {
		k.kindCMDPath = kindPath
		return nil
	}

	kindPath, err = findKindPath()
	if err != nil {
		return err
	}
	k.kindCMDPath = kindPath

	return nil
}

func (k *KindOperator) SetExecCommand(execCommand func(string, ...string) *exec.Cmd) {
	k.execCommand = execCommand
}

func (k *KindOperator) KindVersion() (string, error) {
	b, err := k.execCommand(k.kindCMDPath, "version").CombinedOutput()
	if err != nil {
		return "", err
	}
	klog.V(1).Infof("get kind version %s", b)
	info := strings.Split(string(b), " ")
	// the output of "kind version" is like:
	// kind v0.11.1 go1.17.7 linux/amd64
	ver := info[1]
	return ver, nil
}

func (k *KindOperator) KindLoadDockerImage(out io.Writer, clusterName, image string, nodeNames []string) error {
	nodeArgs := strings.Join(nodeNames, ",")
	klog.V(1).Infof("load image %s to nodes %s in cluster %s", image, nodeArgs, clusterName)
	cmd := k.execCommand(k.kindCMDPath, "load", "docker-image", image, "--name", clusterName, "--nodes", nodeArgs)
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}
	if err := cmd.Run(); err != nil {
		klog.Errorf("failed to load docker image %s to nodes %s in cluster %s, %v", image, nodeArgs, clusterName, err)
		return err
	}
	return nil
}

func (k *KindOperator) KindCreateClusterWithConfig(out io.Writer, configPath string) error {
	cmd := k.execCommand(k.kindCMDPath, "create", "cluster", "--config", configPath, "--kubeconfig", k.kubeconfigPath, "--retain")
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (k *KindOperator) KindDeleteCluster(out io.Writer, name string) error {
	cmd := k.execCommand(k.kindCMDPath, "delete", "cluster", "--name", name)
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkIfKindAt(path string) (bool, string) {
	if p, err := exec.LookPath(path); err == nil {
		return true, p
	}
	return false, ""
}

func findKindPath() (string, error) {

	var kindPath string
	if exist, path := checkIfKindAt("kind"); exist {
		kindPath = path
	}

	if len(kindPath) == 0 {
		return kindPath, fmt.Errorf("cannot find valid kind cmd, try to install it")
	}

	if err := validateKindVersion(kindPath); err != nil {
		return "", err
	}
	return kindPath, nil
}

func validateKindVersion(kindCmdPath string) error {
	tmpOperator := NewKindOperator(kindCmdPath, "")
	ver, err := tmpOperator.KindVersion()
	if err != nil {
		return err
	}
	if !strutil.IsInStringLst(validKindVersions, ver) {
		return fmt.Errorf("invalid kind version: %s, all valid kind versions are: %s",
			ver, strings.Join(validKindVersions, ","))
	}
	return nil
}
