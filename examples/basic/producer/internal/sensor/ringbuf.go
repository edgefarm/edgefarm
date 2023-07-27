package sensor

import (
	"sync"

	"github.com/ci4rail/bogie-pdm/pkg/ringbuf"
)

type samplesRingbuf struct {
	Buf   *ringbuf.Ringbuf[float32]
	mutex *sync.Mutex
}

func (s *samplesRingbuf) Lock() {
	s.mutex.Lock()
}
func (s *samplesRingbuf) Unlock() {
	s.mutex.Unlock()
}
