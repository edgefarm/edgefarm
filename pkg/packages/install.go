/*
Copyright © 2023 EdgeFarm Authors

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

package packages

import (
	"context"
	"fmt"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/s0rg/retry"
	"helm.sh/helm/v3/pkg/repo"

	"k8s.io/klog/v2"

	"github.com/edgefarm/edgefarm/pkg/args"
	"github.com/edgefarm/edgefarm/pkg/k8s"
)

type Spec struct {
	Chart           []*helmclient.ChartSpec
	CreateNamespace bool
	Condition       func() bool
	ValuesFunc      func() string
}

type Helm struct {
	Repo *repo.Entry
	Spec *Spec
}

type Manifest struct {
	Name             string
	Manifest         string
	Condition        func() bool
	WaitForCondition bool
}

type Packages struct {
	Helm     []*Helm
	Manifest []*Manifest
}

func InstallHelmSpec(client helmclient.Client, spec *helmclient.ChartSpec) error {
	// Install helm chart with repeat mechanism if failed to install
	try := retry.New(
		retry.Count(3),
		retry.Sleep(time.Second*2),
		retry.Verbose(true),
	)
	if err := try.Single(fmt.Sprintf("InstallOrUpgradeChart for chart %s", spec.ChartName),
		func() error {
			release, err := client.InstallOrUpgradeChart(context.Background(), spec)
			if err != nil {
				return err
			}
			if release == nil {
				return fmt.Errorf("release failed")
			}
			klog.Infof("chart: %s: installed release %s", spec.ChartName, release.Name)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (h *Helm) Uninstall() error {
	if h.Spec.Condition != nil {
		if !h.Spec.Condition() {
			klog.Info("condition not met, skipping helm chart uninstallation for: ")
			for _, spec := range h.Spec.Chart {
				klog.Infof("chart: %s", spec.ChartName)
			}
			return nil
		}
	}
	for _, spec := range h.Spec.Chart {
		client, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
			Options: &helmclient.Options{
				Namespace: spec.Namespace,
				Debug:     true,
				Linting:   false,
				DebugLog:  klog.Infof,
			},
			RestConfig: args.KubeConfigRestConfig,
		})
		if err != nil {
			return err
		}
		if err := client.UninstallRelease(spec); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manifest) Install() error {
	if m.Condition != nil {
		if m.WaitForCondition {
			try := retry.New(
				retry.Count(30),
				retry.Sleep(time.Second*2),
				retry.Verbose(true),
			)
			if err := try.Single(fmt.Sprintf("Installing manifest %s", m.Name),
				func() error {
					return k8s.Apply(m.Manifest)
				}); err != nil {
				return err
			}
			return nil
		}
	} else {
		return k8s.Apply(m.Manifest)
	}
	return nil
}

func (h *Helm) Install() error {
	if h.Spec.Condition != nil {
		if !h.Spec.Condition() {
			klog.Info("condition not met, skipping helm chart installation for: ")
			for _, spec := range h.Spec.Chart {
				klog.Infof("chart: %s", spec.ChartName)
			}
			return nil
		}
	}
	for _, spec := range h.Spec.Chart {
		client, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
			Options: &helmclient.Options{
				Namespace: spec.Namespace,
				Debug:     true,
				Linting:   false,
				DebugLog:  klog.Infof,
			},
			RestConfig: args.KubeConfigRestConfig,
		})
		if h.Repo != nil {
			if h.Repo.URL != "" {
				if err := client.AddOrUpdateChartRepo(*h.Repo); err != nil {
					return err
				}
			}
		}
		if err != nil {
			return err
		}

		if h.Spec.ValuesFunc != nil {
			spec.ValuesYaml = h.Spec.ValuesFunc()
		}

		if h.Spec.CreateNamespace {
			klog.Infof("chart: %s: creating namespace %s", spec.ChartName, spec.Namespace)
			_, err := k8s.CreateNamespaceIfNotExist(spec.Namespace)
			if err != nil {
				return err
			}
		}
		if err := InstallHelmSpec(client, spec); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packages) Install() error {
	if p.Helm != nil {
		for _, helm := range p.Helm {
			if err := helm.Install(); err != nil {
				return err
			}
		}
	}
	if p.Manifest != nil {
		for _, manifest := range p.Manifest {
			if err := manifest.Install(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Packages) Uninstall() error {
	if p.Helm != nil {
		for _, helm := range p.Helm {
			if err := helm.Uninstall(); err != nil {
				return err
			}
		}
	}
	return nil
}

func InstallBase() error {
	for _, pkg := range Base {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func InstallDependencies() error {
	for _, pkg := range ClusterDependencies {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func Install(packages []Packages) error {
	for _, pkg := range packages {
		if err := pkg.Install(); err != nil {
			return err
		}
	}
	return nil
}

func Uninstall(packages []Packages) error {
	for _, pkg := range packages {
		if err := pkg.Uninstall(); err != nil {
			return err
		}
	}
	return nil
}

func InstallPackage(pkg Packages) error {
	return pkg.Install()
}

func InstallPackageByName(name string, pkgs []Packages) error {
	for _, pkg := range pkgs {
		if pkg.Helm != nil {
			for _, helm := range pkg.Helm {
				if helm.Repo.Name == name {
					return helm.Install()
				}
			}
		}
	}
	return fmt.Errorf("package %s not found", name)
}
