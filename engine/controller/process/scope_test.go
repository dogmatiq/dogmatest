package process_test

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogmatest/engine/controller/process"
	"github.com/dogmatiq/dogmatest/engine/envelope"
	"github.com/dogmatiq/dogmatest/engine/fact"
	"github.com/dogmatiq/enginekit/fixtures"
	handlerkit "github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type eventScope", func() {
	var (
		handler    *fixtures.ProcessMessageHandler
		controller *Controller
		event      = envelope.New(
			fixtures.MessageA1,
			message.EventRole,
		)
	)

	BeforeEach(func() {
		handler = &fixtures.ProcessMessageHandler{
			RouteEventToInstanceFunc: func(
				_ context.Context,
				m dogma.Message,
			) (string, bool, error) {
				switch m.(type) {
				case fixtures.MessageA:
					return "<instance>", true, nil
				default:
					panic(dogma.UnexpectedMessage)
				}
			},
		}
		controller = NewController("<name>", handler)
	})

	When("the instance does not exist", func() {
		Describe("func Root", func() {
			It("panics", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.Root()
					return nil
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						time.Now(),
						event,
					)
				}).To(Panic())
			})
		})

		Describe("func Begin", func() {
			It("returns true", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					Expect(s.Begin()).To(BeTrue())
					return nil
				}

				_, err := controller.Handle(
					context.Background(),
					fact.Ignore,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
			})

			It("records a fact", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.Begin()
					return nil
				}

				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).To(ContainElement(
					fact.ProcessInstanceBegun{
						HandlerName: "<name>",
						InstanceID:  "<instance>",
						Root:        &fixtures.ProcessRoot{},
						Envelope:    event,
					},
				))
			})
		})

		Describe("func End", func() {
			It("panics", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.End()
					return nil
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						time.Now(),
						event,
					)
				}).To(Panic())
			})
		})

		Describe("func ExecuteCommand", func() {
			It("panics", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.ExecuteCommand(fixtures.MessageB1)
					return nil
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						time.Now(),
						event,
					)
				}).To(Panic())
			})
		})

		Describe("func ScheduleTimeout", func() {
			It("panics", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.ScheduleTimeout(fixtures.MessageB1, time.Now())
					return nil
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						time.Now(),
						event,
					)
				}).To(Panic())
			})
		})
	})

	When("the instance exists", func() {
		BeforeEach(func() {
			handler.HandleEventFunc = func(
				_ context.Context,
				s dogma.ProcessEventScope,
				_ dogma.Message,
			) error {
				s.Begin()
				return nil
			}

			_, err := controller.Handle(
				context.Background(),
				fact.Ignore,
				time.Now(),
				envelope.New(
					fixtures.MessageA2, // use a different message to create the instance
					message.EventRole,
				),
			)

			Expect(err).ShouldNot(HaveOccurred())
		})

		Describe("func Root", func() {
			It("returns the root", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					Expect(s.Root()).To(Equal(&fixtures.ProcessRoot{}))
					return nil
				}

				_, err := controller.Handle(
					context.Background(),
					fact.Ignore,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("func Begin", func() {
			It("returns false", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					Expect(s.Begin()).To(BeFalse())
					return nil
				}

				_, err := controller.Handle(
					context.Background(),
					fact.Ignore,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
			})

			It("does not record a fact", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.Begin()
					return nil
				}

				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).NotTo(ContainElement(
					BeAssignableToTypeOf(fact.ProcessInstanceBegun{}),
				))
			})
		})

		Describe("func End", func() {
			It("records a fact", func() {
				handler.HandleEventFunc = func(
					_ context.Context,
					s dogma.ProcessEventScope,
					_ dogma.Message,
				) error {
					s.End()
					return nil
				}

				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).To(ContainElement(
					fact.ProcessInstanceEnded{
						HandlerName: "<name>",
						InstanceID:  "<instance>",
						Root:        &fixtures.ProcessRoot{},
						Envelope:    event,
					},
				))
			})
		})

		Describe("func ExecuteCommand", func() {
			BeforeEach(func() {
				fn := handler.HandleEventFunc
				handler.HandleEventFunc = func(
					ctx context.Context,
					s dogma.ProcessEventScope,
					m dogma.Message,
				) error {
					if err := fn(ctx, s, m); err != nil {
						return err
					}

					s.ExecuteCommand(fixtures.MessageB1)
					return nil
				}
			})

			It("records a fact", func() {
				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).To(ContainElement(
					fact.CommandExecutedByProcess{
						HandlerName: "<name>",
						InstanceID:  "<instance>",
						Root:        &fixtures.ProcessRoot{},
						Envelope:    event,
						CommandEnvelope: event.NewCommand(
							fixtures.MessageB1,
							envelope.Origin{
								HandlerName: "<name>",
								HandlerType: handlerkit.ProcessType,
								InstanceID:  "<instance>",
							},
						),
					},
				))
			})
		})

		Describe("func ScheduleTimeout", func() {
			t := time.Now()

			BeforeEach(func() {
				fn := handler.HandleEventFunc
				handler.HandleEventFunc = func(
					ctx context.Context,
					s dogma.ProcessEventScope,
					m dogma.Message,
				) error {
					if err := fn(ctx, s, m); err != nil {
						return err
					}

					s.ScheduleTimeout(fixtures.MessageB1, t)
					return nil
				}
			})

			It("records a fact", func() {
				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					time.Now(),
					event,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).To(ContainElement(
					fact.TimeoutScheduledByProcess{
						HandlerName: "<name>",
						InstanceID:  "<instance>",
						Root:        &fixtures.ProcessRoot{},
						Envelope:    event,
						TimeoutEnvelope: event.NewTimeout(
							fixtures.MessageB1,
							t,
							envelope.Origin{
								HandlerName: "<name>",
								HandlerType: handlerkit.ProcessType,
								InstanceID:  "<instance>",
							},
						),
					},
				))
			})
		})
	})

	Describe("func InstanceID", func() {
		It("returns the instance ID", func() {
			called := false
			handler.HandleEventFunc = func(
				_ context.Context,
				s dogma.ProcessEventScope,
				_ dogma.Message,
			) error {
				called = true
				Expect(s.InstanceID()).To(Equal("<instance>"))
				return nil
			}

			_, err := controller.Handle(
				context.Background(),
				fact.Ignore,
				time.Now(),
				event,
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(called).To(BeTrue())
		})
	})

	Describe("func Log", func() {
		BeforeEach(func() {
			handler.HandleEventFunc = func(
				_ context.Context,
				s dogma.ProcessEventScope,
				_ dogma.Message,
			) error {
				s.Log("<format>", "<arg-1>", "<arg-2>")
				return nil
			}
		})

		It("records a fact", func() {
			buf := &fact.Buffer{}
			_, err := controller.Handle(
				context.Background(),
				buf,
				time.Now(),
				event,
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.Facts).To(ContainElement(
				fact.MessageLoggedByProcess{
					HandlerName: "<name>",
					InstanceID:  "<instance>",
					Root:        &fixtures.ProcessRoot{},
					Envelope:    event,
					LogFormat:   "<format>",
					LogArguments: []interface{}{
						"<arg-1>",
						"<arg-2>",
					},
				},
			))
		})
	})
})
