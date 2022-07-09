package apps

import (
	"fmt"
	"strings"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = g.Describe("Edge Simple App Deployment", g.Serial, func() {

	g.Describe("Edge Single App", func() {
		var (
			f *fw.Framework
		)
		g.BeforeEach(func() {
			f = fw.DefaultFramework
		})
		g.AfterEach(func() {
			fw.ExpectNoError(f.RemoveNodeLabels(nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName)
			fw.ExpectNoError(f.WaitForNoPodsInNamespace(testingNameSpace, dsPollTimeout))
		})

		g.It("App can be deployed on a specific edge node", func() {
			numNodes := 1
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
			l := getFirstPodLog(testingNameSpace, podNamePrefix1, 10)
			fmt.Printf("Logs received are: %v\n", l)
			Expect(l).To(ContainSubstring("Hello"))
		})

		g.It("App Component version can be changed", func() {
			numNodes := 1
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
			img, err := getFirstPodImage(testingNameSpace, podNamePrefix1)
			fw.ExpectNoError(err)
			Expect(img).To(ContainSubstring("glibc"))

			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple_modversion.yaml")
			err = wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
				img, err = getFirstPodImage(testingNameSpace, podNamePrefix1)
				if err != nil {
					return false, nil
				}
				if strings.Contains(img, "musl") {
					return true, nil
				}
				return false, nil
			})
			fw.ExpectNoError(err)
		})

		g.It("App can be deployed on multiple edge nodes", func() {
			numNodes := 2
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
		})
	})

	g.Describe("Edge Multi App", func() {
		var (
			f *fw.Framework
		)
		g.BeforeEach(func() {
			f = fw.DefaultFramework
		})
		g.AfterEach(func() {
			fw.ExpectNoError(f.RemoveNodeLabels(nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName)
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName2)
			fw.ExpectNoError(f.WaitForNoPodsInNamespace(testingNameSpace, dsPollTimeout))
		})

		g.It("Multiple apps can be deployed on a specific edge node", func() {
			numNodes := 1
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/app2.yaml")
			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
			l := getFirstPodLog(testingNameSpace, podNamePrefix1, 10)
			Expect(l).To(ContainSubstring("Hello"))

			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix2, numNodes))
			l = getFirstPodLog(testingNameSpace, podNamePrefix2, 10)
			Expect(l).To(ContainSubstring("HelloFromApp2-A"))

			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix3, numNodes))
			l = getFirstPodLog(testingNameSpace, podNamePrefix3, 10)
			Expect(l).To(ContainSubstring("HelloFromApp2-B"))
		})

		g.It("Multiple apps can be deployed on a many edge nodes", func() {
			numNodes := 3
			fw.ExpectNoError(f.LabelReadyEdgeNodes(testingNameSpace, numNodes, nodeLabelKey))
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/app2.yaml")
			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
			l := getFirstPodLog(testingNameSpace, podNamePrefix1, 10)
			Expect(l).To(ContainSubstring("Hello"))

			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix2, numNodes))
			l = getFirstPodLog(testingNameSpace, podNamePrefix2, 10)
			Expect(l).To(ContainSubstring("HelloFromApp2-A"))

			fw.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix3, numNodes))
			l = getFirstPodLog(testingNameSpace, podNamePrefix3, 10)
			Expect(l).To(ContainSubstring("HelloFromApp2-B"))
		})
	})
})
