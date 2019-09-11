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
			expectsDescription: "ER_PARSER_TRACE was added in 8.0.2.",
		},
		{
			name:               "Error scheduler killing",
			errString:          "MY-010054",
			expects:            ErrSchedulerKilling,
			expectsString:      "ER_SCHEDULER_KILLING",
			expectsDescription: "ER_SCHEDULER_KILLING was added in 8.0.2.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := FromErrorString(tc.errString)

			if r != tc.expects {
				t.Fatal("Expected other string")
			}

			if r.String() != tc.expectsString {
				t.Fatal("Expected other string")
			}

			if r.Description() != tc.expectsDescription {
				t.Fatal("Expected other description")
			}
		})
	}
}
