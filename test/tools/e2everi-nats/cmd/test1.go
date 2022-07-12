package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edgefarm/edgefarm/test/pkg/msg"
	nats "github.com/nats-io/nats.go"
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
	userCreds     string
	consumerName  string
	stream        string
	t1OkThreshold int
)

type messageEnvelope struct {
	Data msg.Message
	//Id   string
}

func runTest1(cmd *cobra.Command, args []string) {
	address := os.Getenv("NATS_ADDRESS")
	if address == "" {
		address = nats.DefaultURL
	}

	opts := []nats.Option{}

	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	nc, err := nats.Connect(address, opts...)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v\n", err)
	}
	log.Printf("Connected to %s\n", address)
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error creating jetstream context: %v\n", err)
	}
	subjectFilter := ">"

	cinfo, err := js.AddConsumer(stream, &nats.ConsumerConfig{Durable: consumerName, AckPolicy: nats.AckExplicitPolicy, ReplayPolicy: nats.ReplayInstantPolicy})
	if err != nil {
		log.Fatalf("Error creating consumer: %v\n", err)
	}
	log.Printf("consumer created: %v\n", cinfo)

	sub, err := js.PullSubscribe(subjectFilter, consumerName, nats.Bind(stream, consumerName))
	if err != nil {
		log.Fatalf("Error creating subscription: %v\n", err)
	}

	verifier := msg.NewVerifier(t1OkThreshold)

	for {
		messages, err := sub.Fetch(100, nats.PullOpt(nats.MaxWait(time.Second*10)))
		if err == nats.ErrTimeout {
			log.Fatalf("Timeout reading stream: %v\n", err)
		} else if err != nil {
			log.Fatalf("Error reading stream: %v\n", err)
		} else {
			log.Printf("Fetched %d messages\n", len(messages))
			for _, m := range messages {
				var data messageEnvelope
				err := json.Unmarshal(m.Data, &data)
				if err != nil {
					log.Fatalf("Error unmarshal: %v\n", err)
				}

				log.Printf("got message on subject %s: %v\n",  m.Subject, data)
				err = verifier.VerifyMessage(m.Subject, data.Data)
				if err != nil {
					log.Println(err)
				}

				err = m.Ack()
				if err != nil {
					log.Fatalf("Error ack message: %v\n", err)
				}
			}
		}
		dumpProducers(verifier)
	}
}

func dumpProducers(v *msg.Verifier) {
	for k, p := range v.ProducerMap {
		e := ""
		if p.State == msg.FinishedError {
			e = fmt.Sprintf(": %v", p.Err)
		}
		fmt.Printf("\t%s: %s%s\n", k, p.State, e)
	}
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	test1Cmd.Flags().StringVarP(&userCreds, "creds", "", "", "nats credentials")
	test1Cmd.Flags().StringVarP(&consumerName, "consumer", "", "e2everi-nats", "nats consumer name")
	test1Cmd.Flags().StringVarP(&stream, "stream", "s", "e2e-network_export-stream-aggregate", "nats stream name")
	test1Cmd.Flags().IntVarP(&t1OkThreshold, "okthreshold", "o", 1000, "say ok after successfully received this number of messages")
}
