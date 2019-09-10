package mysqlerrnum

import (
	"regexp"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

// ErrorNumber represents the error number from mysql.MySQLError
type ErrorNumber int

// ErrorString represents 'MY-' prefixed errors.
type ErrorString string

var (
	mySQLErrRe = regexp.MustCompile(`Error (\d+):`)
)

// FromError takes an error, tries to cast it as a mysql.MySQLError defined by
// go-sql-driver and return the ErrorNumber corresponding to said number.
func FromError(err error) ErrorNumber {
	if e, ok := err.(*mysql.MySQLError); ok {
		return FromNumber(int(e.Number))
	}

	return FromString(err.Error())
}

// FromString tries to parse a string and get the error number from said string.
func FromString(s string) ErrorNumber {
	nString := mySQLErrRe.FindStringSubmatch(s)

	if len(nString) != 2 {
		return ErrUnknownMySQLError
	}

	if n, err := strconv.Atoi(nString[1]); err == nil {
		return ErrorNumber(n)
	}

	return ErrUnknownMySQLError
}
