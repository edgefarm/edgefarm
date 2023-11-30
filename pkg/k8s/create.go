/*
Copyright 2020 The OpenYurt Authors.

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

package k8s

import (
	"context"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	yaml "github.com/edgefarm/edgefarm/pkg/k8s/yaml"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"
)

const (
// SystemNamespace = "kube-system"
// DefaultWaitServantJobTimeout specifies the timeout value of waiting for the ServantJob to be succeeded
// DefaultWaitServantJobTimeout = time.Minute * 5
)

func processCreateErr(kind string, name string, err error) error {
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			klog.V(4).Infof("[WARNING] %s/%s is already in cluster, skip to prepare it", kind, name)
			return nil
		}
		return fmt.Errorf("fail to create the %s/%s: %w", kind, name, err)
	}
	klog.V(4).Infof("%s/%s is created", kind, name)
	return nil
}

// CreateSecretFromYaml creates the Secret from the yaml template.
func CreateSecretFromYaml(cliSet kubeclientset.Interface, ns, saTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(saTmpl))
	if err != nil {
		return err
	}
	se, ok := obj.(*corev1.Secret)
	if !ok {
		return fmt.Errorf("fail to assert secret: %w", err)
	}
	_, err = cliSet.CoreV1().Secrets(ns).Create(context.Background(), se, metav1.CreateOptions{})
	return processCreateErr("secret", se.Name, err)
}

// CreateServiceAccountFromYaml creates the ServiceAccount from the yaml template.
func CreateServiceAccountFromYaml(cliSet kubeclientset.Interface, ns, saTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(saTmpl))
	if err != nil {
		return err
	}
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return fmt.Errorf("fail to assert serviceaccount: %w", err)
	}
	_, err = cliSet.CoreV1().ServiceAccounts(ns).Create(context.Background(), sa, metav1.CreateOptions{})
	return processCreateErr("serviceaccount", sa.Name, err)
}

// CreateClusterRoleFromYaml creates the ClusterRole from the yaml template.
func CreateClusterRoleFromYaml(cliSet kubeclientset.Interface, crTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(crTmpl))
	if err != nil {
		return err
	}
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return fmt.Errorf("fail to assert clusterrole: %w", err)
	}
	_, err = cliSet.RbacV1().ClusterRoles().Create(context.Background(), cr, metav1.CreateOptions{})
	return processCreateErr("clusterrole", cr.Name, err)
}

// CreateClusterRoleBindingFromYaml creates the ClusterRoleBinding from the yaml template.
func CreateClusterRoleBindingFromYaml(cliSet kubeclientset.Interface, crbTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(crbTmpl))
	if err != nil {
		return err
	}
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert clusterrolebinding: %w", err)
	}
	_, err = cliSet.RbacV1().ClusterRoleBindings().Create(context.Background(), crb, metav1.CreateOptions{})
	return processCreateErr("clusterrolebinding", crb.Name, err)
}

// CreateRoleBindingFromYaml creates the RoleBinding from the yaml template.
func CreateRoleBindingFromYaml(cliSet kubeclientset.Interface, crbTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(crbTmpl))
	if err != nil {
		return err
	}
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return fmt.Errorf("fail to assert rolebinding: %w", err)
	}
	_, err = cliSet.RbacV1().RoleBindings("kube-system").Create(context.Background(), rb, metav1.CreateOptions{})
	return processCreateErr("rolebinding", rb.Name, err)
}

// CreateConfigMapFromYaml creates the ConfigMap from the yaml template.
func CreateConfigMapFromYaml(cliSet kubeclientset.Interface, ns, cmTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(cmTmpl))
	if err != nil {
		return err
	}
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return fmt.Errorf("fail to assert configmap: %w", err)
	}
	_, err = cliSet.CoreV1().ConfigMaps(ns).Create(context.Background(), cm, metav1.CreateOptions{})
	return processCreateErr("configmap", cm.Name, err)
}

// CreateDeployFromYaml creates the Deployment from the yaml template.
func CreateDeployFromYaml(cliSet kubeclientset.Interface, ns, dplyTmpl string, ctx interface{}) error {
	ycmdp, err := tmplutil.SubsituteTemplate(dplyTmpl, ctx)
	if err != nil {
		return err
	}
	dpObj, err := yaml.YamlToObject([]byte(ycmdp))
	if err != nil {
		return err
	}
	dply, ok := dpObj.(*appsv1.Deployment)
	if !ok {
		return errors.New("fail to assert Deployment")
	}
	_, err = cliSet.AppsV1().Deployments(ns).Create(context.Background(), dply, metav1.CreateOptions{})
	return processCreateErr("deployment", dply.Name, err)
}

// CreateServiceFromYaml creates the Service from the yaml template.
func CreateServiceFromYaml(cliSet kubeclientset.Interface, ns, svcTmpl string) error {
	obj, err := yaml.YamlToObject([]byte(svcTmpl))
	if err != nil {
		return err
	}
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return fmt.Errorf("fail to assert service: %w", err)
	}
	_, err = cliSet.CoreV1().Services(ns).Create(context.Background(), svc, metav1.CreateOptions{})
	return processCreateErr("service", svc.Name, err)
}
