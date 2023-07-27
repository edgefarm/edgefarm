// Package daprpubsub is a wrapper around the DAPR client
package daprpubsub

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// Connection is the DAPR connection
type Connection struct {
	client dapr.Client
	nodeID string
}

// New creates a new instance of DaprPubSub
func New(address string, nodeID string) (*Connection, error) {
	var client dapr.Client
	var err error
	for {
		client, err = dapr.NewClientWithAddress(address)
		if err == nil {
			break
		}
		log.Warn().Msgf("dapr client: %s, retrying...", err)
		time.Sleep(2 * time.Second)
	}

	return &Connection{
		client: client,
		nodeID: nodeID,
	}, nil
}

// PubExport publishes a message to the export stream
func (c *Connection) PubExport(subject string, data []byte) error {
	ctx := context.Background()
	subject = fmt.Sprintf("%s.%s", c.nodeID, subject)

	// application/octet-stream is import. Then DAPR will encode payload in base64
	err := c.client.PublishEvent(ctx, "streams", subject, data, dapr.PublishEventWithContentType("application/octet-stream"))
	return err
}
