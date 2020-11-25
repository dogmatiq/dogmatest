package testkit_test

import (
	"context"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/testkit"
	"github.com/dogmatiq/testkit/engine"
	"github.com/dogmatiq/testkit/internal/testingmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ToExecuteCommand()", func() {
	var (
		testingT *testingmock.T
		app      dogma.Application
		test     *Test
	)

	BeforeEach(func() {
		testingT = &testingmock.T{
			FailSilently: true,
		}

		app = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")

				c.RegisterAggregate(&AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", "<aggregate-key>")
						c.ConsumesCommandType(MessageR{}) // R = record an event
						c.ProducesEventType(MessageN{})
					},
					RouteCommandToInstanceFunc: func(dogma.Message) string {
						return "<instance>"
					},
					HandleCommandFunc: func(
						_ dogma.AggregateRoot,
						s dogma.AggregateCommandScope,
						_ dogma.Message,
					) {
						s.RecordEvent(MessageN1)
					},
				})

				c.RegisterProcess(&ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", "<process-key>")
						c.ConsumesEventType(MessageE{})   // E = event
						c.ConsumesEventType(MessageN{})   // N = (do) nothing
						c.ProducesCommandType(MessageC{}) // C = command
					},
					RouteEventToInstanceFunc: func(
						context.Context,
						dogma.Message,
					) (string, bool, error) {
						return "<instance>", true, nil
					},
					HandleEventFunc: func(
						_ context.Context,
						s dogma.ProcessEventScope,
						m dogma.Message,
					) error {
						if _, ok := m.(MessageE); ok {
							s.Begin()
							s.ExecuteCommand(MessageC1)
						}
						return nil
					},
				})
			},
		}
	})

	DescribeTable(
		"expectation behavior",
		func(
			a Action,
			e Expectation,
			ok bool,
			rm reportMatcher,
			options ...TestOption,
		) {
			test = Begin(testingT, app, options...)
			test.Expect(a, e)
			rm(testingT)
			Expect(testingT.Failed()).To(Equal(!ok))
		},
		Entry(
			"command executed as expected",
			RecordEvent(MessageE1),
			ToExecuteCommand(MessageC1),
			expectPass,
			expectReport(
				`✓ execute a specific 'fixtures.MessageC' command`,
			),
		),
		Entry(
			"no matching command executed",
			RecordEvent(MessageE1),
			ToExecuteCommand(MessageX1),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageX' command`,
				``,
				`  | EXPLANATION`,
				`  |     none of the engaged handlers executed the expected command`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify the logic within the '<process>' process message handler`,
			),
		),
		Entry(
			"no messages produced at all",
			RecordEvent(MessageN1),
			ToExecuteCommand(MessageX1),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageX' command`,
				``,
				`  | EXPLANATION`,
				`  |     no messages were produced at all`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify the logic within the '<process>' process message handler`,
			),
		),
		Entry(
			"no commands produced at all",
			ExecuteCommand(MessageR1),
			ToExecuteCommand(MessageC1),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageC' command`,
				``,
				`  | EXPLANATION`,
				`  |     no commands were executed at all`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify the logic within the '<process>' process message handler`,
			),
		),
		Entry(
			"no matching command executed and all relevant handler types disabled",
			ExecuteCommand(MessageR1),
			ToExecuteCommand(MessageX1),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageX' command`,
				``,
				`  | EXPLANATION`,
				`  |     no relevant handler types were enabled`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • enable process handlers using the EnableHandlerType() option`,
			),
			WithUnsafeOperationOptions(
				engine.EnableProcesses(false),
			),
		),
		Entry(
			"similar command executed with a different type",
			RecordEvent(MessageE1),
			ToExecuteCommand(&MessageC1), // note: message type is pointer
			expectFail,
			expectReport(
				`✗ execute a specific '*fixtures.MessageC' command`,
				``,
				`  | EXPLANATION`,
				`  |     a command of a similar type was executed by the '<process>' process message handler`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • check the message type, should it be a pointer?`,
				`  | `,
				`  | MESSAGE DIFF`,
				`  |     [-*-]fixtures.MessageC{`,
				`  |         Value: "C1"`,
				`  |     }`,
			),
		),
		Entry(
			"similar command executed with a different value",
			RecordEvent(MessageE1),
			ToExecuteCommand(MessageC2),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageC' command`,
				``,
				`  | EXPLANATION`,
				`  |     a similar command was executed by the '<process>' process message handler`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • check the content of the message`,
				`  | `,
				`  | MESSAGE DIFF`,
				`  |     fixtures.MessageC{`,
				`  |         Value: "C[-2-]{+1+}"`,
				`  |     }`,
			),
		),
		Entry(
			"expected message recorded as an event rather than executed as a command",
			ExecuteCommand(MessageR1),
			ToExecuteCommand(MessageN1),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageN' command`,
				``,
				`  | EXPLANATION`,
				`  |     the expected message was recorded as an event by the '<aggregate>' aggregate message handler`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify that the '<aggregate>' aggregate message handler intended to record an event of this type`,
				`  |     • verify that ToExecuteCommand() is the correct expectation, did you mean ToRecordEvent()?`,
			),
		),
		Entry(
			"similar message with a different value recorded as an event rather than executed as a command",
			ExecuteCommand(MessageR1),
			ToExecuteCommand(MessageN2),
			expectFail,
			expectReport(
				`✗ execute a specific 'fixtures.MessageN' command`,
				``,
				`  | EXPLANATION`,
				`  |     a similar message was recorded as an event by the '<aggregate>' aggregate message handler`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify that the '<aggregate>' aggregate message handler intended to record an event of this type`,
				`  |     • verify that ToExecuteCommand() is the correct expectation, did you mean ToRecordEvent()?`,
				`  | `,
				`  | MESSAGE DIFF`,
				`  |     fixtures.MessageN{`,
				`  |         Value: "N[-2-]{+1+}"`,
				`  |     }`,
			),
		),
		Entry(
			"similar message with a different type recorded as an event rather than executed as a command",
			ExecuteCommand(MessageR1),
			ToExecuteCommand(&MessageN2), // note: message type is pointer
			expectFail,
			expectReport(
				`✗ execute a specific '*fixtures.MessageN' command`,
				``,
				`  | EXPLANATION`,
				`  |     a message of a similar type was recorded as an event by the '<aggregate>' aggregate message handler`,
				`  | `,
				`  | SUGGESTIONS`,
				`  |     • verify that the '<aggregate>' aggregate message handler intended to record an event of this type`,
				`  |     • verify that ToExecuteCommand() is the correct expectation, did you mean ToRecordEvent()?`,
				`  |     • check the message type, should it be a pointer?`,
				`  | `,
				`  | MESSAGE DIFF`,
				`  |     [-*-]fixtures.MessageN{`,
				`  |         Value: "N[-2-]{+1+}"`,
				`  |     }`,
			),
		),
	)

	It("panics if the message is invalid", func() {
		Expect(func() {
			ToExecuteCommand(nil)
		}).To(PanicWith("ToExecuteCommand(<nil>): message must not be nil"))
	})
})
