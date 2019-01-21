package engine

import (
	"github.com/dogmatiq/dogmatest/compare"
	"github.com/dogmatiq/dogmatest/engine/fact"
	"github.com/dogmatiq/dogmatest/internal/enginekit/handler"
	"github.com/dogmatiq/dogmatest/render"
)

// Option applies optional settings to an engine.
type Option func(*configurer) error

// WithComparator returns an engine option that specifies the comparator to use.
func WithComparator(c compare.Comparator) Option {
	if c == nil {
		panic("comparator must not be nil")
	}

	return func(cfgr *configurer) error {
		cfgr.comparator = c
		return nil
	}
}

// WithRenderer returns an engine option that specifies the renderer to use.
func WithRenderer(r render.Renderer) Option {
	if r == nil {
		panic("renderer must not be nil")
	}

	return func(cfgr *configurer) error {
		cfgr.renderer = r
		return nil
	}
}

// DispatchOption applies optional settings while dispatching a message.
type DispatchOption func(*dispatchOptions) error

// WithObserver returns a dispatch option that registers the given observer
// while the message is being dispatched.
//
// Multiple observers can be registered during a single dispatch.
func WithObserver(o fact.Observer) DispatchOption {
	if o == nil {
		panic("observer must not be nil")
	}

	return func(do *dispatchOptions) error {
		do.observers = append(do.observers, o)
		return nil
	}
}

// EnableHandlerType returns a dispatch option that enables or disables handlers
// of the given type.
//
// By default, aggregates and processes are enabled, and integrations and
// projections are disabled.
func EnableHandlerType(t HandlerType, enable bool) DispatchOption {
	t.MustValidate()

	return func(do *dispatchOptions) error {
		do.enabledHandlers[t] = enable
		return nil
	}
}

// dispatchOptions is a container for the options set via DispatchOption values.
type dispatchOptions struct {
	observers       fact.ObserverSet
	enabledHandlers map[handler.Type]bool
}

// newDispatchOptions returns a new dispatchOptions with the given options.
func newDispatchOptions(options []DispatchOption) (*dispatchOptions, error) {
	do := &dispatchOptions{
		enabledHandlers: map[handler.Type]bool{
			handler.AggregateType:   true,
			handler.ProcessType:     true,
			handler.IntegrationType: false,
			handler.ProjectionType:  false,
		},
	}

	for _, opt := range options {
		if err := opt(do); err != nil {
			return nil, err
		}
	}

	return do, nil
}

// HandlerType is an enumeration of the types of messages handlers.
type HandlerType = handler.Type

const (
	// AggregateType is the handler type for dogma.AggregateMessageHandler.
	AggregateType = handler.AggregateType

	// ProcessType is the handler type for dogma.ProcessMessageHandler.
	ProcessType = handler.ProcessType

	// IntegrationType is the handler type for dogma.IntegrationMessageHandler.
	IntegrationType = handler.IntegrationType

	// ProjectionType is the handler type for dogma.ProjectionMessageHandler.
	ProjectionType = handler.ProjectionType
)
