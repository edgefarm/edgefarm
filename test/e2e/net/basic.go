package net

import (
	"fmt"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"github.com/edgefarm/edgefarm/test/pkg/msg"
	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
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
			err := wait.PollImmediate(time.Second*5, dsPollTimeout, func() (bool, error) {
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
				fmt.Sprintf("%s_%s", net1Name, streamName))
			fw.ExpectNoError(err)
			defer sub.Close()

			taggedNodes, _ := f.GetTaggedNodes(nodeLabelKey)
			nodeName := taggedNodes[0]

			verifier := msg.NewVerifier(1000)
			start := time.Now()
			for {
				m, subject, err := sub.Next(time.Second * 1)
				fw.ExpectNoError(err, "error getting message from nats")
				o.Expect(m).NotTo(o.BeNil(), "no message received")
				verifier.VerifyMessage(subject, *m)

				pub, state := verifier.PublisherStatus(nodeName, net1Name, pubId1)
				if state == msg.FinishedOk {
					break
				} else if state == msg.FinishedError {
					g.Fail(fmt.Sprintf("Publisher verification failed: %v", pub.Err))
				} else if time.Since(start) > dsPollTimeout {
					g.Fail("Publisher verification timed out")
				}
			}
		})
	})
})
