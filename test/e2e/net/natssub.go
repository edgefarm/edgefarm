package net

import (
	"encoding/json"
	"time"

	"github.com/edgefarm/edgefarm/test/pkg/msg"
	nats "github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	nc  *nats.Conn
	sub *nats.Subscription
}

type messageEnvelope struct {
	Data msg.Message
}

// NewNatsSubscriber creates a new NatsSubscriber for 
// subject pattern provided as consumer on stream
func NewNatsSubscriber(natsUrl string, creds string, subject string, consumer string, stream string) (*NatsSubscriber, error) {
	opts := []nats.Option{nats.UserCredentials(creds)}
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	sub, err := js.SubscribeSync(subject, nats.Durable(consumer), nats.MaxDeliver(3), nats.BindStream(stream))
	if err != nil {
		return nil, err
	}
	return &NatsSubscriber{nc, sub}, nil
}

// Close closes the connection to the NATS server.
func (n *NatsSubscriber) Close() {
	n.nc.Close()
}

// Next returns the next message from the subscription.
// If no message is available, it returns nil.
func (n *NatsSubscriber) Next(timeout time.Duration) (*msg.Message, error) {
	m, err := n.sub.NextMsg(timeout)
	if err == nats.ErrTimeout {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var data messageEnvelope
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			return nil, err
		}
		err = m.Ack()
		if err != nil {
			return nil, err
		}
		return &data.Data, nil
	}
}
