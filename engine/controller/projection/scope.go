package projection

import (
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/testkit/engine/envelope"
	"github.com/dogmatiq/testkit/engine/fact"
)

// scope is an implementation of dogma.ProjectionEventScope.
type scope struct {
	config   configkit.RichProjection
	observer fact.Observer
	event    *envelope.Envelope
}

func (s *scope) RecordedAt() time.Time {
	return s.event.CreatedAt
}

func (s *scope) Log(f string, v ...interface{}) {
	s.observer.Notify(fact.MessageLoggedByProjection{
		HandlerName:  s.config.Identity().Name,
		Handler:      s.config.Handler(),
		Envelope:     s.event,
		LogFormat:    f,
		LogArguments: v,
	})
}
