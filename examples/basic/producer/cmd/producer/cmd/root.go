/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"os"

	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/edgefarm/edgefarm/examples/basic/producer/internal/daprpubsub"
	"github.com/edgefarm/edgefarm/examples/basic/producer/internal/export"
	"github.com/edgefarm/edgefarm/examples/basic/producer/internal/sensor"
)

var rootCmd = &cobra.Command{
	Use:   "simulated sensor producer example",
	Short: "Example producer of simulated sensor values",
	Long:  `Example producer of simulated sensor values`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.999Z07:00"})
	pubsub := pubsub.New(100)

	var exporter export.Exporter
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		log.Fatal().Msg("NODE_NAME not set")
	}

	address := os.Getenv("DAPR_GRPC_ADDRESS")
	if address == "" {
		log.Fatal().Msg("DAPR_GRPC_ADDRESS not set")
	}

	exporter, err := daprpubsub.New(address, nodeName)
	if err != nil {
		log.Fatal().Msgf("dapr pubsub: %s", err)
	}

	sensorUnit := sensor.New(pubsub, exporter, &sensor.Configuration{
		SampleRate:     1000,
		RingBufEntries: 1024,
	})
	if sensorUnit != nil {
		sensorUnit.Run("sensor")
	}
	select {}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("Execute Root cmd: %s", err)
	}
}

func init() {
	cobra.OnInitialize()
}
