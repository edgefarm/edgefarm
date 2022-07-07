package apps

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func nodeIsReady(n *corev1.Node) bool {
	for _, c := range n.Status.Conditions {
		if c.Type == "Ready" {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func labelNodes(nameSpace string, numNodes int) error {
	f := framework.DefaultFramework
	nodes, err := f.GetNodes(metav1.ListOptions{})
	if err != nil {
		return err
	}
	i := 0
	for _, n := range nodes.Items {
		if nodeIsReady(&n) {
			_, ok := n.ObjectMeta.Labels[edgeLabelKey]
			if ok {
				err := f.SetNodeLabel(&n, nodeLabelKey, "")
				if err != nil {
					return err
				}
				i++
				if i == numNodes {
					break
				}
			}
		}
	}
	if i < numNodes {
		return fmt.Errorf("cannot tag requested number of nodes")
	}
	return nil
}

func removeNodeLabels() error {
	f := framework.DefaultFramework
	nodes, err := f.GetNodes(metav1.ListOptions{})
	if err != nil {
		return err
	}
	// remove test labels from nodes
	for _, n := range nodes.Items {
		err := f.RemoveNodeLabel(&n, nodeLabelKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func waitPodsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) error {
	err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
		return podsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, expectedPodNamePrefix, expectedInstances)
	})
	return err
}

func podsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) (bool, error) {
	f := framework.DefaultFramework

	podNodes := getRunningPodsNodeNames(nameSpace, expectedPodNamePrefix)

	// get list of tagged nodes
	nods, err := f.GetNodes(metav1.ListOptions{})
	framework.ExpectNoError(err)

	taggedNodes := make([]string, 0)
	for _, n := range nods.Items {
		_, ok := n.ObjectMeta.Labels[labelKey]
		if ok {
			taggedNodes = append(taggedNodes, n.Name)
		}
	}

	// check if the two lists are identical
	sort.Strings(podNodes)
	sort.Strings(taggedNodes)

	fmt.Printf("podNodes: %v, taggedNodes: %v\n", podNodes, taggedNodes)

	if len(taggedNodes) == 0 {
		return false, fmt.Errorf("no tagged nodes")
	}

	if len(podNodes) > len(taggedNodes) {
		return false, fmt.Errorf("too many pods started")
	}

	if len(podNodes) == len(taggedNodes) {
		for i := 0; i < len(podNodes); i++ {
			if podNodes[i] != taggedNodes[i] {
				return false, fmt.Errorf("pod started on wrong node")
			}
		}
		if len(podNodes) != expectedInstances {
			return false, fmt.Errorf("wrong number of instances")
		}
		return true, nil
	}

	return false, nil
}

func waitForNoPodsInNamespace(nameSpace string) error {
	err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
		f := framework.DefaultFramework

		// get all running pods in namespace
		pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		if len(pods.Items) == 0 {
			return true, nil
		}
		return false, nil
	})
	return err
}

func getRunningPodsNodeNames(nameSpace string, podNamePrefix string) []string {
	f := framework.DefaultFramework
	podNodes := make([]string, 0)

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return podNodes
	}

	// get of nodes where the pods are running
	for _, p := range pods.Items {
		if strings.HasPrefix(p.Name, podNamePrefix) && p.Status.Phase == corev1.PodRunning {
			podNodes = append(podNodes, p.Spec.NodeName)
		}
	}
	return podNodes
}

func getRunningPodsNames(nameSpace string, podNamePrefix string) []string {
	f := framework.DefaultFramework
	podNames := make([]string, 0)

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return podNames
	}

	// get of nodes where the pods are running
	for _, p := range pods.Items {
		if strings.HasPrefix(p.Name, podNamePrefix) && p.Status.Phase == corev1.PodRunning {
			podNames = append(podNames, p.Name)
		}
	}
	return podNames
}

func getPodImage(nameSpace string, podName string, containerName string) (string, error) {

	pod, err := getPodByName(nameSpace, podName)
	if err != nil {
		return "", err
	}

	for _, c := range pod.Spec.Containers {
		if c.Name == containerName {
			return c.Image, nil
		}
	}
	return "", fmt.Errorf("containerName not found")
}

func getPodByName(nameSpace string, podName string) (*corev1.Pod, error) {
	f := framework.DefaultFramework
	return f.ClientSet.CoreV1().Pods(nameSpace).Get(f.Context, podName, metav1.GetOptions{})
}

func getFirstPodName(nameSpace string, podNamePrefix string) (string, error) {
	podNames := getRunningPodsNames(nameSpace, podNamePrefix)

	if len(podNames) < 1 {
		return "", fmt.Errorf("no pods found")
	}
	return podNames[0], nil
}

func getFirstPodImage(nameSpace string, podNamePrefix string) (string, error) {
	podName, err := getFirstPodName(nameSpace, podNamePrefix)
	if err != nil {
		return "", err
	}

	return getPodImage(nameSpace, podName, podNamePrefix)
}

func getFirstPodLog(nameSpace string, podNamePrefix string, tailLines int) string {
	pn, err := getFirstPodName(nameSpace, podNamePrefix)
	if err != nil {
		return ""
	}
	p, err := getPodByName(nameSpace, pn)
	if err != nil {
		return ""
	}
	f := framework.DefaultFramework
	l, err := f.GetPodLog(*p, int64(tailLines))
	if err != nil {
		return ""
	}
	return l
}
