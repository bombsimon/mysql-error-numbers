package mysqlerrnum

import (
	"github.com/go-sql-driver/mysql"
)

// ErrorNumber represents the error number from mysql.MySQLError
type ErrorNumber int

// ErrorString represents 'MY-' prefixed errors.
type ErrorString string

// FromError takes an error, tries to cast it as a mysql.MySQLError defined by
// go-sql-driver and return the ErrorNumber corresponding to said number.
func FromError(err error) ErrorNumber {
	if e, ok := err.(*mysql.MySQLError); ok {
		return FromNumber(int(e.Number))
	}

	return ErrorNumber(-1)
}
