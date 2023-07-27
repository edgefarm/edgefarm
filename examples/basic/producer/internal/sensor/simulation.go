package sensor

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/ci4rail/bogie-pdm/pkg/ringbuf"
)

func (s *Unit) simulatedSampling(sampleRate float64, rbEntries int32) (*samplesRingbuf, error) {
	s.logger.Info().Msgf("simulation sampler starting")

	// create ringbuf
	rb := &samplesRingbuf{
		Buf:   ringbuf.New[float32](int(rbEntries)),
		mutex: &sync.Mutex{},
	}

	go func() {
		for {
			// random delay between 800 and 3000ms
			time.Sleep(time.Duration(rand.Intn(4200)+800) * time.Millisecond)
			// generate sine wave with decreasing amplitude
			amplitude := 1.0 * rand.Float64()
			startAmpl := amplitude

			for i := int32(0); i < rbEntries; i++ {
				amplitude -= (1 / float64(rbEntries) * startAmpl)
				value := float32(amplitude * math.Sin(float64(float64(i)*sampleRate/20)*math.Pi/float64(rbEntries)/2))
				rb.Buf.Push(value)
			}
			s.ps.Pub("trigger", "trigger")
		}
	}()
	return rb, nil
}
