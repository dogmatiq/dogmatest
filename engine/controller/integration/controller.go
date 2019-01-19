package integration

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogmatest/engine/controller"
	"github.com/dogmatiq/dogmatest/render"
)

// Controller is an implementation of engine.Controller for
// dogma.IntegrationMessageHandler implementations.
type Controller struct {
	name     string
	handler  dogma.IntegrationMessageHandler
	renderer render.Renderer
}

// NewController returns a new controller for the given handler.
func NewController(
	n string,
	h dogma.IntegrationMessageHandler,
	r render.Renderer,
) *Controller {
	return &Controller{
		name:     n,
		handler:  h,
		renderer: r,
	}
}

// Handle handles a message.
func (c *Controller) Handle(ctx context.Context, s controller.Scope) error {
	panic("not implemented")
}

// Reset clears the state of the controller.
func (c *Controller) Reset() {
}