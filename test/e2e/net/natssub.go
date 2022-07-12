package net

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/edgefarm/edgefarm/test/pkg/msg"
	nats "github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	nc        *nats.Conn
	sub       *nats.Subscription
	credsFile *os.File
}

type messageEnvelope struct {
	Data msg.Message
}

type VeriferFunc func(subject string, m msg.Message) error

// NewNatsSubscriber creates a new NatsSubscriber for
// subject pattern provided as consumer on stream
func NewNatsSubscriber(natsUrl string, creds string, subject string, consumer string, stream string) (*NatsSubscriber, error) {
	file, err := ioutil.TempFile("", "e2e-nats-subscriber")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(file.Name(), []byte(creds), 0600)
	if err != nil {
		return nil, err
	}

	opts := []nats.Option{nats.UserCredentials(file.Name())}
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		os.Remove(file.Name())
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		os.Remove(file.Name())
		nc.Close()
		return nil, err
	}
	_, err = js.AddConsumer(stream, &nats.ConsumerConfig{Durable: consumer, AckPolicy: nats.AckExplicitPolicy, ReplayPolicy: nats.ReplayInstantPolicy})
	if err != nil {
		os.Remove(file.Name())
		nc.Close()
		return nil, fmt.Errorf("can't add consumer: %v", err)
	}
	sub, err := js.PullSubscribe(subject, consumer, nats.Bind(stream, consumer))
	if err != nil {
		os.Remove(file.Name())
		nc.Close()
		return nil, fmt.Errorf("can't create subscription: %v", err)
	}
	return &NatsSubscriber{nc, sub, file}, nil
}

// Close closes the connection to the NATS server.
func (n *NatsSubscriber) Close() {
	n.nc.Close()
	os.Remove(n.credsFile.Name())
}

// NextBatch verifies the next batch of messages.
func (n *NatsSubscriber) NextBatch(nMessages int, timeout time.Duration, verifier VeriferFunc) error {
	messages, err := n.sub.Fetch(nMessages, nats.PullOpt(nats.MaxWait(timeout)))
	if err == nats.ErrTimeout {
		fmt.Print("Timeout reading stream\n")
		return nil
	} else if err != nil {
		return err
	} else {
		for _, m := range messages {
			var data messageEnvelope
			err := json.Unmarshal(m.Data, &data)
			if err != nil {
				return err
			}
			//g.GinkgoWriter.Printf("Got message: %v\n", data.Data)
			err = m.Ack()
			if err != nil {
				return err
			}
			verifier(m.Subject, data.Data)
		}
		return nil
	}
}
