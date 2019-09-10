# MySQL Error Numbers

[![](https://godoc.org/github.com/bombsimon/mysql-error-numbers?status.svg)](http://godoc.org/github.com/bombsimon/mysql-error-numbers)
[![Build Status](https://travis-ci.org/bombsimon/mysql-error-numbers.svg?branch=master)](https://travis-ci.org/bombsimon/mysql-error-numbers)

This is a generated package defining all MySQL errors from the documentation as
constants. The reasone for this is to make the code more readable and
potentially help describe the errors.

The name and number are fetched from version 8.0 of the [MySQL
documentation](https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html).

## Usage

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "github.com/bombsimon/mysql-error-numbers"
)

func main() {
    db, _ := sql.Open("mysql", "user:password@/dbname")

    // Pass the error itself.
    if _, err := db.Query("INSERT INTO... <BAD STATEMENT>"); err != nil {
        if mysqlerrnum.FromError(err) == mysqlerrnum.ErrDupEntry {
            panic("You fool, duplicate!")
        }
    }

    // Or do the conversion.
    _, err := db.Query("INSERT INTO... <BAD STATEMENT>")
    if e, ok := err.(*MySQLError); ok {
        if mysqlerrnum.FromNumber(e.Number) != mysqlerrnum.ErrNoReferencedRow {
            panic("Oops, expected ErNoReferencedRow")
        }
    }
}
```
