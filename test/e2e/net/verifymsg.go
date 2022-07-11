package net

import (
	"fmt"
	"time"

	"github.com/edgefarm/edgefarm/test/pkg/msg"
	g "github.com/onsi/ginkgo/v2"
)

type publisherExpect struct {
	subject string
	id      int
}

func verifyPublishers(sub *NatsSubscriber, expectedPublishers []publisherExpect, expectedMessages int) error {

	verifier := msg.NewVerifier(expectedMessages)
	start := time.Now()
	for {
		m, subject, err := sub.Next(time.Second * 1)
		if err != nil {
			g.GinkgoWriter.Printf("error getting message from nats %v\n", err)
			time.Sleep(time.Second * 1)
			continue
		}

		if m == nil {
			g.GinkgoWriter.Printf("no message received\n")
		} else {
			verifier.VerifyMessage(subject, *m)
		}

		finishCount := int(0)
		for _, p := range expectedPublishers {
			_, state := verifier.PublisherStatus(p.subject, p.id)

			if state == msg.FinishedOk || state == msg.FinishedError {
				finishCount++
			}
		}
		if finishCount == len(expectedPublishers) {
			break
		}
		if time.Since(start) > dsPollTimeout {
			return fmt.Errorf("publisher verification timed out")
		}
	}
	for _, p := range expectedPublishers {
		pub, state := verifier.PublisherStatus(p.subject, p.id)
		if state != msg.FinishedOk {
			return fmt.Errorf("publisher %s %d verification failed: %v", p.subject, p.id, pub.Err)
		}
	}
	return nil
}
