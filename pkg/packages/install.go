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
	"io"
	"net/http"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/s0rg/retry"
	"helm.sh/helm/v3/pkg/repo"

	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

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
	URI              string
	Condition        func() bool
	PreHook          func(manifest string) (string, error)
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

func (h *Helm) Uninstall(kubeconfig *rest.Config) error {
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
			RestConfig: kubeconfig,
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

func (m *Manifest) getManifest() (string, error) {
	klog.Infof("Intalling manifest for %s\n", m.Name)
	if m.URI != "" {
		klog.Infof("Getting from %s\n", m.URI)
		return downloadFromURI(m.URI)
	}
	return m.Manifest, nil
}

func (m *Manifest) Install(kubeconfig *rest.Config) error {
	manifest, err := m.getManifest()
	if err != nil {
		return err
	}
	if m.PreHook != nil {
		manifest, err = m.PreHook(manifest)
		if err != nil {
			return err
		}
	}
	if m.Condition != nil {
		if m.WaitForCondition {
			try := retry.New(
				retry.Count(30),
				retry.Sleep(time.Second*2),
				retry.Verbose(true),
			)
			if err := try.Single(fmt.Sprintf("Installing manifest %s", m.Name),
				func() error {
					if c := m.Condition(); !c {
						return fmt.Errorf("condition not met")
					}
					return nil
				}); err != nil {
				return err
			}
			return k8s.Apply(kubeconfig, manifest)
		}
	}
	return k8s.Apply(kubeconfig, manifest)
}

func downloadFromURI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Status code
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (h *Helm) Install(kubeconfig *rest.Config) error {
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
			RestConfig: kubeconfig,
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
			_, err := k8s.CreateNamespaceIfNotExist(kubeconfig, spec.Namespace)
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

func (p *Packages) Install(kubeconfig *rest.Config) error {
	if p.Helm != nil {
		for _, helm := range p.Helm {
			if err := helm.Install(kubeconfig); err != nil {
				return err
			}
		}
	}
	if p.Manifest != nil {
		for _, manifest := range p.Manifest {
			if err := manifest.Install(kubeconfig); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Packages) Uninstall(kubeconfig *rest.Config) error {
	if p.Helm != nil {
		for _, helm := range p.Helm {
			if err := helm.Uninstall(kubeconfig); err != nil {
				return err
			}
		}
	}
	return nil
}

func Install(kubeconfig *rest.Config, packages []Packages) error {
	for _, pkg := range packages {
		if err := pkg.Install(kubeconfig); err != nil {
			return err
		}
	}
	return nil
}

func Uninstall(kubeconfig *rest.Config, packages []Packages) error {
	for _, pkg := range packages {
		if err := pkg.Uninstall(kubeconfig); err != nil {
			return err
		}
	}
	return nil
}
