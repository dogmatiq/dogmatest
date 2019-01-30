package fact

import (
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogmatest/engine/envelope"
	"github.com/dogmatiq/enginekit/handler"
)

// UnroutableMessageDispatched indicates that Engine.Dispatch() has been called
// with a message that is not routed to any handlers.
//
// Note that when dispatch is called with an unroutable message, it is unknown
// whether it was intended to be a command or an event.
type UnroutableMessageDispatched struct {
	// Message is the message that was dispatched.
	Message         dogma.Message
	Now             time.Time
	EnabledHandlers map[handler.Type]bool
}

// MessageDispatchBegun indicates that Engine.Dispatch() has been called with a
// message that is able to be routed to at least one handler.
type MessageDispatchBegun struct {
	Envelope        *envelope.Envelope
	Now             time.Time
	EnabledHandlers map[handler.Type]bool
}

// MessageDispatchCompleted indicates that a call Engine.Dispatch() has completed.
type MessageDispatchCompleted struct {
	Envelope        *envelope.Envelope
	Now             time.Time
	Error           error
	EnabledHandlers map[handler.Type]bool
}

// MessageHandlingBegun indicates that a message is about to be handled by a
// specific handler.
type MessageHandlingBegun struct {
	HandlerName string
	HandlerType handler.Type
	Envelope    *envelope.Envelope
}

// MessageHandlingCompleted indicates that a message has been handled by a
// specific handler, either successfully or unsucessfully.
type MessageHandlingCompleted struct {
	HandlerName string
	HandlerType handler.Type
	Envelope    *envelope.Envelope
	Error       error
}

// MessageHandlingSkipped indicates that a message has been not been handled by
// a specific handler, because handlers of that type are disabled.
type MessageHandlingSkipped struct {
	HandlerName string
	HandlerType handler.Type
	Envelope    *envelope.Envelope
}
