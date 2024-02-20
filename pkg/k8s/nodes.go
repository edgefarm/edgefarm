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

package k8s

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"regexp"
	"time"

	"github.com/s0rg/retry"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	yaml "sigs.k8s.io/yaml"
)

var (
	DefaultEdgeFarmEdgeWorkerLabel = metav1.LabelSelector{
		MatchLabels: map[string]string{
			"openyurt.io/is-edge-worker": "true",
		},
	}

	DefaultEdgeNodeTaint = v1.Taint{
		Key:    "edgefarm.io",
		Value:  "",
		Effect: v1.TaintEffectNoSchedule,
	}
)

func DeleteNodepool(kubeconfig *rest.Config, name string) error {
	dynamic, err := GetDynamicClient(kubeconfig)
	if err != nil {
		return err
	}
	return dynamic.Resource(schema.GroupVersionResource{
		Group:    "apps.openyurt.io",
		Version:  "v1beta1",
		Resource: "nodepools",
	}).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func DeleteNode(kubeconfig *rest.Config, name string) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}
	return clientset.CoreV1().Nodes().Delete(context.Background(), name, metav1.DeleteOptions{})
}

// GetAllNodes returns a slice with all nodes
func GetAllNodes(kubeconfig *rest.Config) ([]v1.Node, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

// GetNodes returns a slice of nodes matching the given selector.
func GetNodes(kubeconfig *rest.Config, selector *metav1.LabelSelector) ([]v1.Node, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(selector),
	})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

func NodeExists(kubeconfig *rest.Config, name string) (bool, error) {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return false, err

	}
	_, err = clientset.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ValidatePhysicalNodeName validates the given name not to be anything like regex `edgefarm-control-plane` or `edgefarm-worker.*`
func ValidatePhysicalNodeName(name string) error {
	if name == "edgefarm-control-plane" {
		return fmt.Errorf("cannot delete node '%s'", name)
	}

	re, err := regexp.Compile(`edgefarm-worker.*`)
	if err != nil {
		return err
	}

	if re.MatchString(name) {
		return fmt.Errorf("cannot delete node '%s'", name)
	}
	return nil
}

func GetEdgeNodes(kubeconfig *rest.Config) ([]v1.Node, error) {
	return GetNodes(kubeconfig, &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"openyurt.io/is-edge-worker": "true",
		},
	})
}

func GetCloudNodes(kubeconfig *rest.Config) ([]v1.Node, error) {
	return GetNodes(kubeconfig, &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"openyurt.io/is-edge-worker": "false",
		},
	})
}

func CheckNodeTaint(node v1.Node, taint v1.Taint) bool {
	if node.Spec.Taints == nil {
		return false
	}

	for _, t := range node.Spec.Taints {
		if t.Key == taint.Key && t.Value == taint.Value && t.Effect == taint.Effect {
			return true
		}
	}
	return false
}

func AnnotateNodes(kubeconfig *rest.Config, nodes []v1.Node, annotations map[string]string) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.Annotations == nil {
			node.Annotations = map[string]string{}
		}
		for k, v := range annotations {
			node.Annotations[k] = v
		}
		try := retry.New(
			retry.Count(30),
			retry.Sleep(time.Second),
			retry.Verbose(true),
		)
		if err := try.Single("Updating node", func() error {
			fresh, err := clientset.CoreV1().Nodes().Get(context.Background(), node.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			fresh.Annotations = node.Annotations
			if _, err := clientset.CoreV1().Nodes().Update(context.Background(), fresh, metav1.UpdateOptions{}); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

	}
	return nil
}

func LabelNodes(kubeconfig *rest.Config, nodes []v1.Node, labels map[string]string) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.Labels == nil {
			node.Labels = map[string]string{}
		}
		for k, v := range labels {
			node.Labels[k] = v
		}
		try := retry.New(
			retry.Count(30),
			retry.Sleep(time.Second),
			retry.Verbose(true),
		)
		if err := try.Single("Updating node", func() error {
			fresh, err := clientset.CoreV1().Nodes().Get(context.Background(), node.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			fresh.Labels = node.Labels
			if _, err := clientset.CoreV1().Nodes().Update(context.Background(), fresh, metav1.UpdateOptions{}); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

	}
	return nil
}

func TaintNodes(kubeconfig *rest.Config, nodes []v1.Node, taint v1.Taint) error {
	clientset, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.Spec.Taints == nil {
			node.Spec.Taints = []v1.Taint{}
		}
		if CheckNodeTaint(node, taint) {
			continue
		}
		try := retry.New(
			retry.Count(30),
			retry.Sleep(time.Second),
			retry.Verbose(true),
		)
		if err := try.Single("Updating node", func() error {
			fresh, err := clientset.CoreV1().Nodes().Get(context.Background(), node.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			fresh.Spec.Taints = append(fresh.Spec.Taints, taint)
			if _, err := clientset.CoreV1().Nodes().Update(context.Background(), fresh, metav1.UpdateOptions{}); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

const (
	nodePoolTemplate = `apiVersion: apps.openyurt.io/v1beta1
kind: NodePool
metadata: 
  labels: 
    monitor.edgefarm.io/metrics: default
    openyurt.io/node-pool-type: edge
  name: {{.Name}}
spec: 
  selector: 
  matchLabels: 
    apps.openyurt.io/nodepool: {{.Name}}
  type: Edge`
)

// HandleNodePool creates a nodepool resource and create the corresponding label on the node
func HandleNodePool(kubeconfig *rest.Config, node v1.Node) error {
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	client, err := GetClientset(kubeconfig)
	if err != nil {
		return err
	}

	try := retry.New(
		retry.Count(30),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single("Updating node", func() error {
		fresh, err := client.CoreV1().Nodes().Get(context.Background(), node.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		fresh.Labels["apps.openyurt.io/desired-nodepool"] = fresh.Name
		if _, err := client.CoreV1().Nodes().Update(context.Background(), fresh, metav1.UpdateOptions{}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	// handle nodepool template
	type Values struct {
		Name string
	}
	values := Values{Name: node.Name}
	tmpl, err := template.New("test").Parse(nodePoolTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return err
	}
	j, err := yaml.YAMLToJSON(buf.Bytes())
	if err != nil {
		return err
	}

	manifest := &unstructured.Unstructured{}
	if err := manifest.UnmarshalJSON(j); err != nil {
		return err
	}

	dynamic, err := GetDynamicClient(nil)
	if err != nil {
		return err
	}

	// only create nodepool if it does not exist
	if _, err := dynamic.Resource(schema.GroupVersionResource{
		Group:    "apps.openyurt.io",
		Version:  "v1beta1",
		Resource: "nodepools",
	}).Get(context.Background(), node.Name, metav1.GetOptions{}); err == nil {
		return nil
	}

	if _, err := dynamic.Resource(schema.GroupVersionResource{
		Group:    "apps.openyurt.io",
		Version:  "v1beta1",
		Resource: "nodepools",
	}).Create(context.Background(), manifest, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil

}

func PrepareEdgeNodes(kubeconfig *rest.Config) error {
	nodes, err := GetEdgeNodes(kubeconfig)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
	}
	klog.Infof("Prepare edge nodes")
	err = AnnotateNodes(kubeconfig, nodes, map[string]string{
		"apps.openyurt.io/binding": "true",
	})
	if err != nil {
		return err
	}

	err = LabelNodes(kubeconfig, nodes, map[string]string{
		"node.edgefarm.io/to-be-converted": "true",
	})
	if err != nil {
		return err
	}

	err = TaintNodes(kubeconfig, nodes, DefaultEdgeNodeTaint)
	if err != nil {
		return err
	}

	return nil
}

func CreateEdgeNodepools(kubeconfig *rest.Config) error {
	nodes, err := GetEdgeNodes(kubeconfig)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err = HandleNodePool(kubeconfig, node)
		if err != nil {
			return err
		}
	}
	return nil
}
