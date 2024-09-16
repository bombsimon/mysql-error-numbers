package mysqlerrnum

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestFromError(t *testing.T) {
	cases := []struct {
		name               string
		err                error
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
		{
			name: "Extract from wrapped error",
			err: fmt.Errorf("wrapped error %w",
				&mysql.MySQLError{
					Number: 1216,
				}),
			expects:            ErrNoReferencedRow,
			expectsString:      "ER_NO_REFERENCED_ROW",
			expectsDescription: "InnoDB reports this error when you try to add a row but there is no parent row, and a foreign key constraint fails. Add the parent row first.",
		},
		{
			name:               "Should not fallback to string",
			err:                errors.New("Error 1216: This should not be catched"),
			expects:            ErrUnknownMySQLError,
			expectsString:      "ER_UNKNOWN_MYSQL_ERROR",
			expectsDescription: "Unknown MySQL error",
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

func TestFromErrorOrString(t *testing.T) {
	cases := []struct {
		name               string
		err                error
		expects            ErrorNumber
		expectsString      string
		expectsDescription string
	}{
		{
			name:               "nil should not panic",
			err:                nil,
			expects:            ErrUnknownMySQLError,
			expectsString:      "ER_UNKNOWN_MYSQL_ERROR",
			expectsDescription: "Unknown MySQL error",
		},
		{
			name:               "non MySQL error returns unknown",
			err:                errors.New("this is not MySQL"),
			expects:            ErrUnknownMySQLError,
			expectsString:      "ER_UNKNOWN_MYSQL_ERROR",
			expectsDescription: "Unknown MySQL error",
		},
		{
			name: "Should succeed with first error",
			err: &mysql.MySQLError{
				Number: 1216,
			},
			expects:            ErrNoReferencedRow,
			expectsString:      "ER_NO_REFERENCED_ROW",
			expectsDescription: "InnoDB reports this error when you try to add a row but there is no parent row, and a foreign key constraint fails. Add the parent row first.",
		},
		{
			name:               "Should succeed with error containing the string",
			err:                errors.New("Error 1216: This should be found"),
			expects:            ErrNoReferencedRow,
			expectsString:      "ER_NO_REFERENCED_ROW",
			expectsDescription: "InnoDB reports this error when you try to add a row but there is no parent row, and a foreign key constraint fails. Add the parent row first.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := FromErrorOrString(tc.err)

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
