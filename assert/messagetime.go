package assert

import (
	"fmt"
	"time"

	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/testkit/render"

	"github.com/dogmatiq/testkit/engine/envelope"
)

type messageTimeCriteria struct {
	StartTime *time.Time
	EndTime   *time.Time
	Ok        bool

	match   *envelope.Envelope
	failure string
}

func (c *messageTimeCriteria) SetMatch(env *envelope.Envelope) bool {
	c.match = env

	if c.StartTime != nil && env.CreatedAt.Before(*c.StartTime) {
		return false
	}

	if c.EndTime != nil && env.CreatedAt.After(*c.EndTime) {
		return false
	}

	return true
}

func (c *messageTimeCriteria) UpdateCriteria(rep *Report) {
	if c.StartTime != nil && c.EndTime != nil {
		rep.Criteria += fmt.Sprintf(
			" between %s and %s",
			c.StartTime.Format(time.RFC3339Nano),
			c.EndTime.Format(time.RFC3339Nano),
		)
	} else if c.StartTime != nil {
		rep.Criteria += fmt.Sprintf(
			" after %s",
			c.StartTime.Format(time.RFC3339Nano),
		)
	} else if c.EndTime != nil {
		rep.Criteria += fmt.Sprintf(
			" before %s",
			c.EndTime.Format(time.RFC3339Nano),
		)
	}
}

// BuildResult builds the assertion result when there is an exact match message
// available but the assertion still failed.
func (a *MessageAssertion) BuildResult(
	match *envelope.Envelope,
	tracker *tracker,
	r render.Renderer,
	rep *Report,
) {
	s := rep.Section(suggestionsSection)

	rep.Explanation = inflect(
		a.role,
		"the <message> was <produced> at %s, which is %s",
		match.Origin.HandlerName,
		match.Origin.HandlerType,
		failure,
	)

	for n, t := range tracker.engaged {
		if t == handler.ProcessType {
			s.AppendListItem("verify the logic within the '%s' %s message handler", n, t)
		}
	}
}
