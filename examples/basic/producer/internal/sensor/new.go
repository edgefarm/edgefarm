package sensor

import (
	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/edgefarm/edgefarm/examples/basic/producer/internal/export"
)

type Configuration struct {
	SampleRate     float64
	RingBufEntries int32
}

type sampler struct {
	rb *samplesRingbuf
}

// Unit is the instance of the SensorUnit
type Unit struct {
	cfg     *Configuration
	logger  zerolog.Logger
	ps      *pubsub.PubSub
	sampler *sampler
	export  export.Exporter
}

// New creates a new instance of the MetricsUnit
func New(ps *pubsub.PubSub, export export.Exporter, cfg *Configuration) *Unit {

	t := &Unit{
		cfg:    cfg,
		ps:     ps,
		logger: log.With().Str("component", "sensor").Logger(),
		export: export,
	}

	return t
}

// Run starts the sensor unit
// It starts the sampler and the publisher with the given subject to publish to
// It returns as soon as all go routines are started
func (s *Unit) Run(subject string) error {
	s.logger.Info().Msg("sensorunit starting")

	// start sensor sampling
	var rb *samplesRingbuf
	// start simulation
	rb, err := s.simulatedSampling(s.cfg.SampleRate, s.cfg.RingBufEntries)
	if err != nil {
		return err
	}

	s.sampler = &sampler{
		rb: rb,
	}
	// start publisher
	s.logger.Info().Msg("about to start publisher")
	s.publisher(subject)

	return nil
}
