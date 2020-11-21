package testkit

import (
	"time"

	"github.com/dogmatiq/testkit/engine"
)

// TestOption applies optional settings to a test.
type TestOption func(*testOptions)

// StartVirtualClockAt returns a test option that sets the initial time of the
// test's virtual clock.
//
// By default, the current system time is used.
func StartVirtualClockAt(t time.Time) TestOption {
	return func(to *testOptions) {
		to.time = t
	}
}

// testOptions is a container for the options set via TestOption values.
type testOptions struct {
	engineOptions    []engine.Option
	operationOptions []engine.OperationOption
	time             time.Time
}

// newTestOptions returns a new testOptions with the given options.
func newTestOptions(options []TestOption) *testOptions {
	ro := &testOptions{
		engineOptions: []engine.Option{
			engine.EnableProjectionCompactionDuringHandling(true),
		},
		operationOptions: []engine.OperationOption{
			engine.EnableIntegrations(false),
			engine.EnableProjections(false),
		},
		time: time.Now(),
	}

	for _, opt := range options {
		opt(ro)
	}

	return ro
}

// WithEngineOptions returns a TestOption that applies optional settings to
// the engine used by the test-runner.
func WithEngineOptions(options ...engine.Option) TestOption {
	return func(to *testOptions) {
		to.engineOptions = append(to.engineOptions, options...)
	}
}

// WithOperationOptions returns a TestOption that applies optional per-operation
// settings when performing assertions.
func WithOperationOptions(options ...engine.OperationOption) TestOption {
	return func(to *testOptions) {
		to.operationOptions = append(to.operationOptions, options...)
	}
}
