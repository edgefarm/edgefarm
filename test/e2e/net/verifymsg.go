package net

import (
	"fmt"
	"time"

	"github.com/edgefarm/edgefarm/test/pkg/msg"
	g "github.com/onsi/ginkgo/v2"
)



func verifyPublishers(sub *NatsSubscriber, expectedProducers []msg.ExpectedProducer, expectedMessages int) error {

	verifier := msg.NewVerifier(expectedProducers, expectedMessages)
	start := time.Now()
	for {
		err := sub.NextBatch(100, time.Second*5, verifier.VerifyMessage)
		if err != nil {
			g.GinkgoWriter.Printf("error getting message batch %v\n", err)
			time.Sleep(time.Second * 1)
			continue
		}
		finishCount := int(0)
		for _, p := range expectedProducers {
			_, state := verifier.PublisherStatus(p.Subject, p.Id)

			if state == msg.FinishedOk || state == msg.FinishedError {
				finishCount++
			}
		}
		if finishCount == len(expectedProducers) {
			break
		}
		if time.Since(start) > dsPollTimeout {
			return fmt.Errorf("publisher verification timed out")
		}
	}
	for _, p := range expectedProducers {
		pub, state := verifier.PublisherStatus(p.Subject, p.Id)
		if state != msg.FinishedOk {
			return fmt.Errorf("publisher %s %d verification failed: %v", p.Subject, p.Id, pub.Err)
		}
	}
	return nil
}
