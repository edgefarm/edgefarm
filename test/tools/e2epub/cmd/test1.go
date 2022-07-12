package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"log"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/spf13/cobra"
)

var test1Cmd = &cobra.Command{
	Use:   "test1",
	Short: "Periodic publish",
	Long: `Periodic publish
	`,
	Run: runTest1,
}

var (
	t1Delay        int
	t1Id           int
	t1Network      string
	t1Subject      string
	t1PrefixNodeID bool
)

type message struct {
	NodeID      string
	UnixSeconds int
	Id          int
	Seq         int
}

func runTest1(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	address := os.Getenv("DAPR_GRPC_ADDRESS")
	if address == "" {
		address = "localhost:3500"
	}

	nodeID := os.Getenv("NODE_NAME")
	if nodeID == "" {
		panic("NODE_NAME not set")
	}

	client, err := dapr.NewClientWithAddress(address)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	subject := t1Subject
	if t1PrefixNodeID {
		subject = fmt.Sprintf("%s.EXPORT.%s", nodeID, t1Subject)
	}
	time.Sleep(time.Second * 25)
	log.Printf("publishing to subject %s\n", subject)

	seq := int(0)
	for {
		m := message{
			NodeID:      nodeID,
			UnixSeconds: int(time.Now().Unix()),
			Id:          t1Id,
			Seq:         seq,
		}
		data, err := json.Marshal(m)

		if err != nil {
			log.Printf("ERR: json encode: %v\n", err)
		} else {
			err := client.PublishEvent(ctx, t1Network, subject, data, dapr.PublishEventWithContentType("application/json"))
			if err != nil {
				log.Printf("ERR: publish: %v\n", err)
			} else {
				log.Printf("published %v\n", string(data))
				seq++
			}
		}
		if seq > 5000 {
			time.Sleep(time.Hour * 5000)
		}

		time.Sleep(time.Millisecond * time.Duration(t1Delay))
	}
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	test1Cmd.Flags().IntVarP(&t1Delay, "delay", "d", 50, "delay between messages in ms")
	test1Cmd.Flags().IntVarP(&t1Id, "id", "i", 1, "id to write into message")
	test1Cmd.Flags().StringVarP(&t1Network, "network", "n", "data-export-network", "network name (dapr component name)")
	test1Cmd.Flags().StringVarP(&t1Subject, "subject", "s", "foo.bar", "publish subject name")
	test1Cmd.Flags().BoolVarP(&t1PrefixNodeID, "prefix", "p", true, "prefix subject with node id")
}
