package net

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = g.SynchronizedBeforeSuite(func() []byte {
	f := fw.DefaultFramework
	f.CreateTestNamespace(testingNameSpace) // ignore error if already exists
	return []byte{}
}, func(d []byte) {
})

var _ = g.Describe("Edgefarm.Network Basic", g.Serial, func() {

	g.Describe("Edge Single NetApp", func() {
		var (
			f *fw.Framework
		)
		g.BeforeEach(func() {
			f = fw.DefaultFramework
		})
		g.AfterEach(func() {
			fw.ExpectNoError(f.RemoveNodeLabels(nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", app1Name)
			fw.ExpectNoError(f.WaitForNoPodsInNamespace(testingNameSpace, dsPollTimeout))
		})

		g.It("Single network app can publish data to main nats", func() {
			numNodes := 1
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))

			// randomize the publisher ID
			pubID := rand.Intn(1000)
			fmt.Printf("using publisher ID %d\n", pubID)

			manifest, err := makePublisherManifest("net/manifest/app1.yaml", "test1", net1Name, 10, pubID)
			fw.ExpectNoError(err)

			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "net/manifest/net1.yaml")
			fw.RunKubectlOrDieInput(kubeConfig, testingNameSpace, manifest, "apply", "-f", "-")

			waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, comp1Name, numNodes)

			sub, err := startupNatsSubscriber(testingNameSpace, app1Name, comp1Name, net1Name)
			fw.ExpectNoError(err)
			defer sub.Close()

			taggedNodes, _ := f.GetTaggedNodes(nodeLabelKey)
			nodeName := taggedNodes[0]
			subject := fmt.Sprintf("EXPORT.%s.foo.bar", nodeName)

			err = verifyPublishers(sub, []publisherExpect{{subject, pubID}}, 1000)
			fw.ExpectNoError(err)
		})
	})
})

// read yaml file, patch it with the provided parameters and return the contents as a string
func makePublisherManifest(yamlName string, testName string, network string, delay int, pubID int) (string, error) {
	f, err := os.ReadFile(yamlName)
	if err != nil {
		return "", err
	}
	var raw map[string]interface{}

	// Unmarshal our input YAML file into empty interface
	if err := yaml.Unmarshal(f, &raw); err != nil {
		return "", err
	}
	spec := raw["spec"].(map[string]interface{})
	components := spec["components"].([]interface{})
	publisher := components[0].(map[string]interface{})
	properties := publisher["properties"].(map[string]interface{})

	// Update our publisher properties
	properties["args"] = []string{testName, fmt.Sprintf("-n=%s", network),
		fmt.Sprintf("-d=%d", delay), fmt.Sprintf("-i=%d", pubID)}

	// Marshal our updated interface into YAML
	b, err := yaml.Marshal(raw)
	fmt.Printf("MANIFEST %s\n", string(b))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func waitPodsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) error {
	err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
		return podsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, expectedPodNamePrefix, expectedInstances)
	})
	return err
}

func podsAreAppliedToAllSelectedNodes(nameSpace string, labelKey string, expectedPodNamePrefix string, expectedInstances int) (bool, error) {
	f := fw.DefaultFramework

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
