package k8s

import (
	"bytes"
	"context"
	"html/template"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

// GetNodes returns a slice of nodes matching the given selector.
func GetNodes(selector metav1.LabelSelector) ([]v1.Node, error) {
	clientset, err := GetClientset(nil)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&selector),
	})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

func GetEdgeNodes() ([]v1.Node, error) {
	return GetNodes(metav1.LabelSelector{
		MatchLabels: map[string]string{
			"openyurt.io/is-edge-worker": "true",
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

func TaintNodes(nodes []v1.Node, taint v1.Taint) error {
	clientset, err := GetClientset(nil)
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
		node.Spec.Taints = append(node.Spec.Taints, taint)
		if _, err := clientset.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{}); err != nil {
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
func HandleNodePool(node v1.Node) error {
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	node.Labels["apps.openyurt.io/desired-nodepool"] = node.Name
	client, err := GetClientset(nil)
	if err != nil {
		return err
	}
	if _, err := client.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{}); err != nil {
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
	if _, err := dynamic.Resource(schema.GroupVersionResource{
		Group:    "apps.openyurt.io",
		Version:  "v1beta1",
		Resource: "nodepools",
	}).Create(context.Background(), manifest, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil

}

func PrepareEdgeNodes() error {
	nodes, err := GetEdgeNodes()
	if err != nil {
		return err
	}
	err = TaintNodes(nodes, DefaultEdgeNodeTaint)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err = HandleNodePool(node)
		if err != nil {
			return err
		}
	}
	return nil
}
