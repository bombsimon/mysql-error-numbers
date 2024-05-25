package mysqlerrnum

import (
	"testing"
)

func TestFromErrorString(t *testing.T) {
	cases := []struct {
		name               string
		errString          string
		expects            ErrorString
		expectsString      string
		expectsDescription string
	}{
		{
			name:               "Error parser trace",
			errString:          "MY-010000",
			expects:            ErrParserTrace,
			expectsString:      "ER_PARSER_TRACE",
			expectsDescription: "Message: Parser saw: %s",
		},
		{
			name:               "Error scheduler killing",
			errString:          "MY-010054",
			expects:            ErrSchedulerKilling,
			expectsString:      "ER_SCHEDULER_KILLING",
			expectsDescription: "Message: Event Scheduler: Killing the scheduler thread, thread id %u",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := FromErrorString(tc.errString)

			if r != tc.expects {
				t.Fatalf("Expected '%s', got '%s'", tc.expects, r)
			}

			if r.String() != tc.expectsString {
				t.Fatalf("Expected '%s', got '%s'", tc.expectsString, r.String())
			}

			if r.Description() != tc.expectsDescription {
				t.Fatalf("Expected '%s', got '%s'", tc.expectsDescription, r.Description())
			}
		})
	}
}
