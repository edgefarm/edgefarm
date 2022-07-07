package apps

import (
	"fmt"
	"strings"
	"time"

	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = ginkgo.Describe("Edge Simple App Deployment", func() {
	ginkgo.AfterEach(func() {
		framework.ExpectNoError(removeNodeLabels())
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "delete", "application", appName)
		framework.ExpectNoError(waitForNoPodsInNamespace(testingNameSpace))
	})

	ginkgo.It("App can be deployed on a specific edge node", func() {
		numNodes := 1
		framework.ExpectNoError(labelNodes(testingNameSpace, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
		l := getFirstPodLog(testingNameSpace, podNamePrefix1, 10)
		fmt.Printf("Logs received are: %v\n", l)
		Expect(l).To(ContainSubstring("Hello"))
	})

	ginkgo.It("App Component version can be changed", func() {
		numNodes := 1
		framework.ExpectNoError(labelNodes(testingNameSpace, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
		img, err := getFirstPodImage(testingNameSpace, podNamePrefix1)
		framework.ExpectNoError(err)
		Expect(img).To(ContainSubstring("glibc"))

		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple_modversion.yaml")
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
		framework.ExpectNoError(err)
	})

	ginkgo.It("App can be deployed on multiple edge nodes", func() {
		numNodes := 2
		framework.ExpectNoError(labelNodes(testingNameSpace, numNodes))
		framework.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "apps/simple.yaml")
		framework.ExpectNoError(waitPodsAreAppliedToAllSelectedNodes(testingNameSpace, nodeLabelKey, podNamePrefix1, numNodes))
	})
})
