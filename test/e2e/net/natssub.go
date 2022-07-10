package net

import (
	"encoding/json"
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
		return nil, err
	}
	sub, err := js.SubscribeSync(subject, nats.Durable(consumer), nats.MaxDeliver(3), nats.BindStream(stream))
	if err != nil {
		os.Remove(file.Name())
		return nil, err
	}
	return &NatsSubscriber{nc, sub, file}, nil
}

// Close closes the connection to the NATS server.
func (n *NatsSubscriber) Close() {
	n.nc.Close()
	os.Remove(n.credsFile.Name())
}

// Next returns the next message from the subscription.
// If no message is available, it returns nil.
// returns the message payload, subject, error
func (n *NatsSubscriber) Next(timeout time.Duration) (*msg.Message, string, error) {
	m, err := n.sub.NextMsg(timeout)
	if err == nats.ErrTimeout {
		return nil, "", nil
	} else if err != nil {
		return nil, "", err
	} else {
		var data messageEnvelope
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			return nil, "", err
		}
		err = m.Ack()
		if err != nil {
			return nil, "", err
		}
		return &data.Data, m.Subject, nil
	}
}
