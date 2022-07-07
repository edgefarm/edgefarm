package apps

import (
	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Edge Multi App Deployment", func() {
	ginkgo.AfterEach(func() {
		framework.ExpectNoError(removeNodeLabels())
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName)
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName2)
		framework.ExpectNoError(waitForNoPodsInNamespace(testingNameSpace))
	})

	ginkgo.It("Multiple apps can be deployed on a specific edge node", func() {
		numNodes := 1
		framework.ExpectNoError(labelNodes(testingNameSpace, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/app2.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
		l := getFirstPodLog(testingNameSpace, podNamePrefix1, 10)
		Expect(l).To(ContainSubstring("Hello"))

		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix2, numNodes))
		l = getFirstPodLog(testingNameSpace, podNamePrefix2, 10)
		Expect(l).To(ContainSubstring("HelloFromApp2-A"))

		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix3, numNodes))
		l = getFirstPodLog(testingNameSpace, podNamePrefix3, 10)
		Expect(l).To(ContainSubstring("HelloFromApp2-B"))
	})
})
