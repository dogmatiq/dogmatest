package aggregate_test

import (
	"context"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogmatest/engine/controller/aggregate"
	"github.com/dogmatiq/dogmatest/engine/envelope"
	"github.com/dogmatiq/dogmatest/engine/fact"
	"github.com/dogmatiq/dogmatest/internal/enginekit/fixtures"
	"github.com/dogmatiq/dogmatest/internal/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type commandScope", func() {
	var (
		handler    *fixtures.AggregateMessageHandler
		controller *Controller
		command    = envelope.New(
			fixtures.MessageA1,
			message.CommandRole,
		)
	)

	BeforeEach(func() {
		handler = &fixtures.AggregateMessageHandler{
			RouteCommandToInstanceFunc: func(m dogma.Message) string {
				switch m.(type) {
				case fixtures.MessageA:
					return "<instance>"
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
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Root()
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						command,
					)
				}).To(Panic())
			})
		})

		Describe("func RecordEvent", func() {
			It("panics", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.RecordEvent(fixtures.MessageB1)
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						command,
					)
				}).To(Panic())
			})
		})

		Describe("func Destroy", func() {
			It("panics", func() {
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					_ dogma.Message,
				) {
					s.Destroy()
				}

				Expect(func() {
					controller.Handle(
						context.Background(),
						fact.Ignore,
						command,
					)
				}).To(Panic())
			})
		})
	})

	When("the instance exists", func() {
		BeforeEach(func() {
			handler.HandleCommandFunc = func(
				s dogma.AggregateCommandScope,
				_ dogma.Message,
			) {
				s.Create()
			}
		})

		Describe("func RecordEvent", func() {
			BeforeEach(func() {
				fn := handler.HandleCommandFunc
				handler.HandleCommandFunc = func(
					s dogma.AggregateCommandScope,
					m dogma.Message,
				) {
					fn(s, m)
					s.RecordEvent(fixtures.MessageB1)
				}
			})

			It("records a fact", func() {
				buf := &fact.Buffer{}
				_, err := controller.Handle(
					context.Background(),
					buf,
					command,
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf.Facts).To(ContainElement(
					fact.EventRecordedByAggregate{
						HandlerName: "<name>",
						InstanceID:  "<instance>",
						Root:        &fixtures.AggregateRoot{},
						Envelope:    command,
						EventEnvelope: command.NewEvent(
							fixtures.MessageB1,
						),
					},
				))
			})
		})
	})

	Describe("func Log", func() {
		BeforeEach(func() {
			handler.HandleCommandFunc = func(
				s dogma.AggregateCommandScope,
				_ dogma.Message,
			) {
				s.Log("<format>", "<arg-1>", "<arg-2>")
			}
		})

		It("records a fact", func() {
			buf := &fact.Buffer{}
			_, err := controller.Handle(
				context.Background(),
				buf,
				command,
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.Facts).To(ContainElement(
				fact.MessageLoggedByAggregate{
					HandlerName: "<name>",
					InstanceID:  "<instance>",
					Root:        &fixtures.AggregateRoot{},
					Envelope:    command,
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
