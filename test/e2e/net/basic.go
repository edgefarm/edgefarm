package net

import (
	"fmt"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"

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

			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "net/manifest/net1.yaml")
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "net/manifest/app1.yaml")

			var creds string
			var natsUrl string
			err := wait.PollImmediate(time.Second, dsPollTimeout, func() (bool, error) {
				var err error
				creds, natsUrl, err = getNatsParams(testingNameSpace, app1Name, comp1Name, net1Name)
				if err != nil {
					g.GinkgoWriter.Printf("Error getting nats params: %v\n", err)
					return false, nil
				}
				return true, nil
			})
			fw.ExpectNoError(err)
			fmt.Printf("Nats params: %s %s\n", creds, natsUrl)

			sub, err := NewNatsSubscriber(natsUrl, creds, "EXPORT.>", "e2e-consumer",
				fmt.Sprintf("%s.%s", net1Name, streamName))
			fw.ExpectNoError(err)
			defer sub.Close()
			
		})
	})
})
