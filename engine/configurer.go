package engine

import (
	"context"

	"github.com/dogmatiq/dogmatest/engine/controller"
	"github.com/dogmatiq/dogmatest/engine/controller/aggregate"
	"github.com/dogmatiq/dogmatest/engine/controller/integration"
	"github.com/dogmatiq/dogmatest/engine/controller/process"
	"github.com/dogmatiq/dogmatest/engine/controller/projection"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/message"
)

type configurer struct {
	engine *Engine
}

func (c *configurer) VisitApplicationConfig(ctx context.Context, cfg *config.ApplicationConfig) error {
	for _, h := range cfg.Handlers {
		if err := h.Accept(ctx, c); err != nil {
			return err
		}
	}

	return nil
}

func (c *configurer) VisitAggregateConfig(_ context.Context, cfg *config.AggregateConfig) error {
	c.registerController(
		aggregate.NewController(
			cfg.HandlerName,
			cfg.Handler,
		),
		message.CommandRole,
		cfg.MessageTypes,
	)

	return nil
}

func (c *configurer) VisitProcessConfig(_ context.Context, cfg *config.ProcessConfig) error {
	c.registerController(
		process.NewController(
			cfg.HandlerName,
			cfg.Handler,
		),
		message.EventRole,
		cfg.MessageTypes,
	)

	return nil
}

func (c *configurer) VisitIntegrationConfig(_ context.Context, cfg *config.IntegrationConfig) error {
	c.registerController(
		integration.NewController(
			cfg.HandlerName,
			cfg.Handler,
		),
		message.CommandRole,
		cfg.MessageTypes,
	)

	return nil
}

func (c *configurer) VisitProjectionConfig(_ context.Context, cfg *config.ProjectionConfig) error {
	c.registerController(
		projection.NewController(
			cfg.HandlerName,
			cfg.Handler,
		),
		message.EventRole,
		cfg.MessageTypes,
	)

	return nil
}

func (c *configurer) registerController(
	ctrl controller.Controller,
	r message.Role,
	types message.TypeSet,
) {
	c.engine.controllers[ctrl.Name()] = ctrl

	for t := range types {
		c.engine.routes[t] = append(c.engine.routes[t], ctrl)
		c.engine.roles[t] = r
	}
}
