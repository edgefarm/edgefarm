package net

import (
	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"
	"github.com/edgefarm/edgefarm/test/pkg/msg"
	. "github.com/onsi/gomega"
)

var _ = g.Describe("Edgefarm.Network Basic", g.Serial, func() {
	g.BeforeSuite(func() {
		f := fw.DefaultFramework
		f.CreateTestNamespace(testingNameSpace) // ignore error if already exists
	})

	g.Describe("Edge Single App", func() {
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

			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "net/manifest/net.yaml")
			fw.RunKubectlOrDie(kubeConfig, testingNameSpace, "apply", "-f", "net/manifest/app.yaml")
			verifier := msg.NewVerifier(1000)
			sub := NewNatsSubscriber()

		})
	})
})
