package mysqlerrnum

import (
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestFromError(t *testing.T) {
	cases := []struct {
		name               string
		err                *mysql.MySQLError
		expects            ErrorNumber
		expectsString      string
		expectsDescription string
	}{
		{
			name: "Error No",
			err: &mysql.MySQLError{
				Number: 1002,
			},
			expects:            ErrNo,
			expectsString:      "ER_NO",
			expectsDescription: "Used in the construction of other messages.",
		},
		{
			name: "Foreign key error",
			err: &mysql.MySQLError{
				Number: 1216,
			},
			expects:            ErrNoReferencedRow,
			expectsString:      "ER_NO_REFERENCED_ROW",
			expectsDescription: "InnoDB reports this error when you try to add a row but there is no parent row, and a foreign key constraint fails. Add the parent row first.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := FromError(tc.err)

			if r != tc.expects {
				t.Fatal("Expected other number")
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

func TestFromString(t *testing.T) {
	cases := []struct {
		name               string
		errString          string
		expects            ErrorNumber
		expectsString      string
		expectsDescription string
	}{
		{
			name:               "Error No",
			errString:          "Error 1002: Some message",
			expects:            ErrNo,
			expectsString:      "ER_NO",
			expectsDescription: "Used in the construction of other messages.",
		},
		{
			name:               "Foreign key error",
			errString:          "Error 1216: Some message",
			expects:            ErrNoReferencedRow,
			expectsString:      "ER_NO_REFERENCED_ROW",
			expectsDescription: "InnoDB reports this error when you try to add a row but there is no parent row, and a foreign key constraint fails. Add the parent row first.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := FromString(tc.errString)

			if r != tc.expects {
				t.Fatal("Expected other number")
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
