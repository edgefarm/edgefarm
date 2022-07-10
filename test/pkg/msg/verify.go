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
	// key "<nodeId>-<subject>-<id>"
	ProducerMap map[string]*Producer
	okThreshold int
}

type Producer struct {
	State            VerifierState
	MessagesReceived int
	Seq              int
	Subject          string
	LastTs           int
	Err              error
}

func NewVerifier(okThreshold int) *Verifier {
	return &Verifier{
		ProducerMap: map[string]*Producer{},
		okThreshold: okThreshold,
	}
}

// VerifyMessage checks m against the states of the known producers
func (v *Verifier) VerifyMessage(subject string, m Message) error {
	key := keyName(m.NodeID, subject, m.Id)

	p, ok := v.ProducerMap[key]

	if !ok {
		// add a new producer
		v.ProducerMap[key] = &Producer{
			MessagesReceived: 0,
			Seq:              m.Seq,
			Subject:          subject,
			LastTs:           m.UnixSeconds,
			State:            Verifying,
		}
		log.Printf("Adding new producer %s for verification. First seq %d\n", key, m.Seq)
	} else {
		if p.State != Verifying {
			// either in error state or verification already ok. Ignore message
			return nil
		} else {
			err := error(nil)
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

func (v *Verifier) PublisherStatus(nodeID string, subject string, id int) (*Producer, VerifierState)  {
	key := keyName(nodeID, subject, id)
	state := Unknown

	p, ok := v.ProducerMap[key]

	if ok {
		state = p.State
	}
	return p, state
}


func keyName(nodeID string, subject string, id int) string {
	return fmt.Sprintf("%s-%s-%d", nodeID, subject, id)
}
