package projection

import (
	"github.com/dogmatiq/dogmatest/engine/envelope"
	"github.com/dogmatiq/dogmatest/engine/fact"
)

// eventScope is an implementation of dogma.ProjectionEventScope.
type eventScope struct {
	name      string
	observers fact.ObserverSet
	event     *envelope.Envelope
}

func (s *eventScope) Log(f string, v ...interface{}) {
	s.observers.Notify(fact.MessageLoggedByProjection{
		HandlerName:  s.name,
		Envelope:     s.event,
		LogFormat:    f,
		LogArguments: v,
	})
}
