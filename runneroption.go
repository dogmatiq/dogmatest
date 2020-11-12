package testkit

import (
	"github.com/dogmatiq/testkit/engine"
)

// RunnerOption applies optional settings to a test runner.
type RunnerOption func(*runnerOptions)

// runnerOptions is a container for the options set via RunnerOption values.
type runnerOptions struct {
	engineOptions []engine.Option
}

// newRunnerOptions returns a new runnerOptions with the given options.
func newRunnerOptions(options []RunnerOption) *runnerOptions {
	ro := &runnerOptions{
		engineOptions: []engine.Option{
			engine.EnableProjectionCompactionDuringHandling(true),
		},
	}

	for _, opt := range options {
		opt(ro)
	}

	return ro
}

// WithEngineOptions returns a RunnerOption that applies optional settings to
// the engine used by the test-runner.
func WithEngineOptions(options ...engine.Option) RunnerOption {
	return func(ro *runnerOptions) {
		ro.engineOptions = append(ro.engineOptions, options...)
	}
}
