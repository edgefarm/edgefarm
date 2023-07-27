package sensor

import (
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/edgefarm/edgefarm/examples/basic/producer/proto/go/sensor/v1"
)

func (s *Unit) publisher(subject string) {
	go func() {
		s.logger.Info().Msgf("publisher running")

		triggerCh := s.ps.Sub("trigger")

		for tr := range triggerCh {
			// received trigger message, publish data
			o := s.publishData()
			s.logger.Debug().Msgf("received trigger %v, published %d bytes", tr, len(o))
			err := s.export.PubExport(subject, o)
			if err != nil {
				s.logger.Error().Msgf("can't publish %v", err)
			}
		}
	}()
}

func (s *Unit) publishData() []byte {
	ts := time.Now()

	s.sampler.rb.Lock()

	nSamples := s.sampler.rb.Buf.Len()
	samples := make([]float32, nSamples)

	for i := 0; i < nSamples; i++ {
		samples[i] = s.sampler.rb.Buf.At(i)
	}
	s.sampler.rb.Unlock()

	b := &pb.Samples{
		TriggerTimestamp: timestamppb.New(ts),
		SampleRate:       s.cfg.SampleRate,
		Samples:          samples,
	}

	out, err := proto.Marshal(b)
	if err != nil {
		s.logger.Error().Msgf("can't marshall %v", err)
	}
	return out
}
