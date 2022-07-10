package apps

import (
	"fmt"
	"sort"
	"time"

	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"k8s.io/apimachinery/pkg/util/wait"
)

func waitPodsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) error {
	err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
		return podsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, expectedPodNamePrefix, expectedInstances)
	})
	return err
}

func podsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) (bool, error) {
	f := framework.DefaultFramework

	podNodes := f.GetRunningPodsNodeNames(nameSpace, expectedPodNamePrefix)

	taggedNodes, err := f.GetTaggedNodes(labelKey)
	if err != nil {
		return false, err
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

func getFirstPodName(nameSpace string, podNamePrefix string) (string, error) {
	f := framework.DefaultFramework
	podNames := f.GetRunningPodsNames(nameSpace, podNamePrefix)

	if len(podNames) < 1 {
		return "", fmt.Errorf("no pods found")
	}
	return podNames[0], nil
}

func getFirstPodImage(nameSpace string, podNamePrefix string) (string, error) {
	f := framework.DefaultFramework
	podName, err := getFirstPodName(nameSpace, podNamePrefix)
	if err != nil {
		return "", err
	}

	return f.GetPodImage(nameSpace, podName, podNamePrefix)
}

func getFirstPodLog(nameSpace string, podNamePrefix string, tailLines int) string {
	f := framework.DefaultFramework
	pn, err := getFirstPodName(nameSpace, podNamePrefix)
	if err != nil {
		return ""
	}
	p, err := f.GetPodByName(nameSpace, pn)
	if err != nil {
		return ""
	}
	l, err := f.GetPodLog(*p, int64(tailLines))
	if err != nil {
		return ""
	}
	return l
}
