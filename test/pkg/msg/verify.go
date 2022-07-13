package msg

import (
	"fmt"
	"log"
)

type VerifierState int

const (
	Unknown VerifierState = iota
	Verifying
	FinishedOk
	FinishedError
)

func (s VerifierState) String() string {
	switch s {
	case Verifying:
		return "verifying"
	case FinishedOk:
		return "ok"
	case FinishedError:
		return "error"
	}
	return "unknown"
}

type Verifier struct {
	// ProducerMap is a map with the detected producers together with their state
	// key "<subject>-<id>"
	ProducerMap       map[string]*Producer
	expectedProducers []ExpectedProducer
	okThreshold       int
}

type ExpectedProducer struct {
	Subject string
	Id      int
}

type Producer struct {
	State            VerifierState
	MessagesReceived int
	Seq              int
	Subject          string
	LastTs           int
	Err              error
}

func NewVerifier(expectedProducers []ExpectedProducer, okThreshold int) *Verifier {
	log.Printf("New Message verifier. expecting producers\n")

	for _, p := range expectedProducers {
		log.Printf(" %s %d\n", p.Subject, p.Id)
	}

	return &Verifier{
		ProducerMap:       map[string]*Producer{},
		expectedProducers: expectedProducers,
		okThreshold:       okThreshold,
	}
}

// VerifyMessage checks m against the states of the known producers
func (v *Verifier) VerifyMessage(subject string, m Message) error {
	key := keyName(subject, m.Id)
	err := error(nil)
	var p *Producer

	if v.isExpectedProducer(subject, m.Id) {
		var ok bool
		p, ok = v.ProducerMap[key]

		if !ok {
			// add a new producer
			v.ProducerMap[key] = &Producer{
				MessagesReceived: 0,
				Seq:              -1,
				Subject:          subject,
				LastTs:           m.UnixSeconds,
				State:            Verifying,
			}
			log.Printf("Adding new producer %s for verification. First seq %d\n", key, m.Seq)
			p = v.ProducerMap[key]
		}
		if p.State != Verifying {
			// either in error state or verification already ok. Ignore message
			return nil
		} else {
			p.MessagesReceived++
			p.Seq++
			p.LastTs = m.UnixSeconds

			if p.Seq != m.Seq {
				err = fmt.Errorf("producer %s: Expected seq %d, got %d", key, p.Seq, m.Seq)
			}
			if err == nil && p.Subject != subject {
				err = fmt.Errorf("producer %s: Expected subject %s, got %s", key, p.Subject, subject)
			}
			if err == nil && p.MessagesReceived >= v.okThreshold {
				p.State = FinishedOk
				log.Printf("### Verification for producer %s ok. Got %d correct messages\n", key, p.MessagesReceived)
			}
			if err != nil {
				p.State = FinishedError
				p.Err = err
				return err
			}
		}
	}
	return nil
}

func (v *Verifier) PublisherStatus(subject string, id int) (*Producer, VerifierState) {
	key := keyName(subject, id)
	state := Unknown

	p, ok := v.ProducerMap[key]

	if ok {
		state = p.State
	}
	return p, state
}

func keyName(subject string, id int) string {
	return fmt.Sprintf("%s-%d", subject, id)
}

func (v *Verifier) isExpectedProducer(subject string, id int) bool {
	for _, p := range v.expectedProducers {
		if p.Subject == subject && p.Id == id {
			return true
		}
	}
	return false
}
