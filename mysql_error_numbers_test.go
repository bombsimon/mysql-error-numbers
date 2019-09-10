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
			expects:            ErNo,
			expectsString:      "ER_NO",
			expectsDescription: "Used in the construction of other messages.",
		},
		{
			name: "Foreign key error",
			err: &mysql.MySQLError{
				Number: 1216,
			},
			expects:            ErNoReferencedRow,
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
