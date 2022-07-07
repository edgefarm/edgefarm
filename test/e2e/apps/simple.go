package apps

import (
	"fmt"
	"sort"
	"time"

	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	nodeLabelKey     = "simple-app" // must match tag in manifest
	appName          = "simple-app" // must match name in manifest
	edgeLabelKey     = "node-role.kubernetes.io/edge"
	dsPollTimeout    = time.Minute * 5
	testingNameSpace = "default" // must match name in manifest
	kubeConfig       = ""
)

var _ = ginkgo.Describe("Edge App Deployment", func() {
	var (
		f     *framework.Framework
		nodes *corev1.NodeList
	)
	ginkgo.JustBeforeEach(func() {
		// use default framework
		f = framework.DefaultFramework
		var err error
		nodes, err = f.GetNodes(metav1.ListOptions{})
		framework.ExpectNoError(err)
	})
	ginkgo.AfterEach(func() {
		// remove test labels from nodes
		for _, n := range nodes.Items {
			framework.ExpectNoError(f.RemoveNodeLabel(&n, nodeLabelKey))
		}
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName)
		framework.ExpectNoError(waitForNoPodsInNamespace(testingNameSpace))
	})

	ginkgo.It("Component can be deployed on a specific edge node", func() {
		numNodes := 1
		framework.ExpectNoError(labelNodes(testingNameSpace, nodes, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodes, nodeLabelKey, numNodes))
	})

	ginkgo.It("Component can be deployed on multiple edge nodes", func() {
		numNodes := 2
		framework.ExpectNoError(labelNodes(testingNameSpace, nodes, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodes, nodeLabelKey, numNodes))
	})
})

func labelNodes(nameSpace string, nodes *corev1.NodeList, numNodes int) error {
	f := framework.DefaultFramework
	i := 0
	for _, n := range nodes.Items {
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
	if i < numNodes {
		return fmt.Errorf("cannot tag requested number of nodes")
	}
	return nil
}

func waitPodsAreAppliedToAllSelectedNodes(nameSpace string, nodes *corev1.NodeList, labelKey string, expectedInstances int) error {
	err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
		return podsAreAppliedToAllSelectedNodes(testingNameSpace, nodes, nodeLabelKey, expectedInstances)
	})
	return err
}

func podsAreAppliedToAllSelectedNodes(nameSpace string, nodes *corev1.NodeList, labelKey string, expectedInstances int) (bool, error) {
	f := framework.DefaultFramework

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	// get of nodes where the pods are running
	podNodes := make([]string, 0)
	for _, p := range pods.Items {
		if p.Status.Phase == corev1.PodRunning {
			podNodes = append(podNodes, p.Spec.NodeName)
		}
	}

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
