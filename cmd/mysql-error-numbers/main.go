package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/iancoleman/strcase"
)

// URLs to the MySQL documentation for the error numbers.
const (
	MySQL57URL = "https://dev.mysql.com/doc/refman/5.7/en/server-error-reference.html"
	MySQL80URL = "https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html"
)

// Literal representsa a found literal from the MySQL documentation.
type Literal struct {
	Number       int
	NumberString string
	Symbol       string
	Constant     string
	State        string
	Extra        string
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(MySQL80URL)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err.Error())
		}
	}()

	if response.StatusCode != http.StatusOK {
		panic(response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Store integer errors and string errors in separate lists. Integer errors
	// are those only containing digits while string errors is prefixed with
	// 'MY-'. The string errors are currently not used.
	var (
		numberLiterals []Literal
		stringLiterals []Literal
		seenLiterals   = map[string]struct{}{}
	)

	// Always add the unknown error.
	numberLiterals = append(numberLiterals, Literal{
		Number:       -1,
		NumberString: "-1",
		Symbol:       "ER_UNKNOWN_MYSQL_ERROR",
		Constant:     "ErrUnknownMySQLError",
		Extra:        "Unknown MySQL error",
	})

	// Fetch every item in the list from '.itemizedlist'.
	document.Find(".itemizedlist li").Each(func(index int, element *goquery.Selection) {
		literal := Literal{}

		// Fetch each literal which is the class name of the code block used to
		// represent the error number/string, the symbol (MySQL error name) and
		// the SQLSTATE.
		element.Find(".literal").Each(func(index int, element *goquery.Selection) {
			switch index {
			case 0:
				number := element.Text()
				n, _ := strconv.Atoi(number)

				literal.Number = n
				literal.NumberString = number
			case 1:
				text := element.Text()

				literal.Symbol = text
				literal.Constant = strcase.ToCamel(
					strings.ToLower(
						// Use 'ERR' instead of 'ER' to follow Go naming
						// convention.
						strings.Replace(text, "ER", "ERR", 1),
					),
				)
			case 2:
				literal.State = element.Text()
			default:
				// 😸
			}
		})

		// Fetch the second p tag in the list to get the additional information
		// about the error.
		element.Find("p").Each(func(index int, element *goquery.Selection) {
			if index == 0 {
				return
			}

			// Normalize and remove spaces and line breaks.
			var (
				text    = strings.TrimSpace(element.Text())
				fields  = strings.Fields(text)
				extra   = strings.Join(fields, " ")
				escaped = strings.Replace(extra, "\"", "\\\"", -1)
			)

			literal.Extra = escaped
		})

		// Ensure we only save the literal once. According to the documentation
		// MySQL re-uses some numbers for legacy reasons, ex.
		// 1120 - ER_WRONG_OUTER_JOIN
		// 1120 - ER_WRONG_OUTER_JOIN_UNUSED
		if _, ok := seenLiterals[literal.NumberString]; ok {
			return
		}

		// Add this literal to seen ones.
		seenLiterals[literal.NumberString] = struct{}{}

		// Skip potential caught li items which isn't a proper error
		// documentation li.
		if literal.NumberString == "" {
			return
		}

		switch literal.Number {
		case 0:
			stringLiterals = append(stringLiterals, literal)
		default:
			numberLiterals = append(numberLiterals, literal)
		}
	})

	args := map[string]interface{}{
		"number_literals": numberLiterals,
		"string_literals": stringLiterals, // Not in use (yet)
		"now":             time.Now().Format("2006-01-02 15:04:05"),
	}

	createFile("number_constants.tmpl", "mysql_error_numbers.gen.go", args)
	createFile("string_constants.tmpl", "mysql_error_strings.gen.go", args)
	createFile("mysql_constants.tmpl", "mysql_error_numbers_constants.gen.go", args)
}

func createFile(templateFile, filename string, data map[string]interface{}) {
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("./cmd/mysql-error-numbers/%s", templateFile)))
	buf := bytes.Buffer{}

	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	fileBytes, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile(filename, fileBytes, 0644)
	fmt.Printf("Generated file: %s\n", filename)
}
