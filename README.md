# MySQL Error Numbers

[![](https://godoc.org/github.com/bombsimon/mysql-error-numbers?status.svg)](http://godoc.org/github.com/bombsimon/mysql-error-numbers)
[![Go](https://github.com/bombsimon/mysql-error-numbers/actions/workflows/go.yml/badge.svg)](https://github.com/bombsimon/mysql-error-numbers/actions/workflows/go.yml)

This is a generated package defining all MySQL errors from the documentation as
constants. The reasone for this is to make the code more readable and
potentially help describe the errors.

The name and number are fetched from version 8.0 of the [MySQL
documentation](https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html).

## Motivation

I had a hard time finding a good way to handle errors and often used a loca map
with numbers to internal error types. I also found
[github.com/VividCortex/mysqlerr](https://github.com/VividCortex/mysqlerr) which
has the error number constants but is not generated (and not maintained?). I
ended up doing this package and would love feedback and PRs to make it a
re-usable way to handle MySQL errors with Go.

## Usage

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "github.com/bombsimon/mysql-error-numbers/v2"
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
