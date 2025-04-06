package clickhouse

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

// ClickZeroTime default clickhouse time
var ClickZeroTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

var bindNamedRe = regexp.MustCompile(`@\w+`) // \w тоже самое, что [A-Za-z0-9_]

type Statement struct {
	Statement     string
	Arg           interface{}
	argCanBeEmpty bool
}

func (s Statement) getNamedArgName() string {
	result := bindNamedRe.Find([]byte(s.Statement))
	return strings.TrimLeft(string(result), "@")
}

func Delimiter(prefix, delimiter string, stmts []Statement) (statementsRow string, args []interface{}) {
	var statements []string
	addArg := func(stmt Statement, arg interface{}) {
		argName := stmt.getNamedArgName()
		if argName != "" {
			args = append(args, clickhouse.Named(argName, arg))
		} else {
			args = append(args, arg)
		}
	}
	deleteEmptyElemInSlice := func(s []string) []string {
		result := make([]string, 0, len(s))
		for _, elem := range s {
			if elem == "" {
				continue
			}
			result = append(result, elem)
		}
		return result
	}

	for _, stmt := range stmts {
		t := reflect.TypeOf(stmt.Arg)
		if t == nil && stmt.Statement != "" {
			statements = append(statements, stmt.Statement)
			continue
		}
		// TODO: Fix linter issue
		switch t.Kind() { //nolint:exhaustive
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(stmt.Arg)
			if stmt.argCanBeEmpty || s.Len() != 0 {
				switch x := stmt.Arg.(type) {
				case []uuid.UUID:
					statements = append(statements, stmt.Statement)
					// addArg(stmt, UUIDArray(x))
					addArg(stmt, x)

				case uuid.UUID:
					// if x != uuid.Nil { // skip empty uuid
					statements = append(statements, stmt.Statement)
					addArg(stmt, x)
					// }

				case []string:
					if !stmt.argCanBeEmpty {
						x = deleteEmptyElemInSlice(x)
					}
					if len(x) != 0 {
						statements = append(statements, stmt.Statement)
						addArg(stmt, x)
					}

				default:
					statements = append(statements, stmt.Statement)
					addArg(stmt, x)
				}
			}

		case reflect.String:
			if stmt.argCanBeEmpty || stmt.Arg.(string) != "" {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}

		case reflect.Bool:
			if stmt.argCanBeEmpty || stmt.Arg.(bool) {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}

		case reflect.Int:
			if stmt.argCanBeEmpty || stmt.Arg.(int) != 0 {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}
		case reflect.Int8:
			if stmt.argCanBeEmpty || stmt.Arg.(int8) != 0 {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}
		case reflect.Int16:
			if stmt.argCanBeEmpty || stmt.Arg.(int16) != 0 {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}
		case reflect.Int32:
			if stmt.argCanBeEmpty || stmt.Arg.(int32) != 0 {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}
		case reflect.Int64:
			if stmt.argCanBeEmpty || stmt.Arg.(int64) != 0 {
				statements = append(statements, stmt.Statement)
				addArg(stmt, stmt.Arg)
			}

		default:
			switch x := stmt.Arg.(type) {
			case time.Time:
				if stmt.argCanBeEmpty || !x.IsZero() {
					statements = append(statements, stmt.Statement)
					addArg(stmt, x)
				}

			default:
				panic(fmt.Sprintf("(click statementsRow) unexpected arg type %T", stmt.Arg))
			}
		}
	}

	if len(statements) != 0 {
		if prefix != "" {
			statementsRow += prefix + " "
		}
		statementsRow += strings.Join(statements, " "+delimiter+" ")
	}

	return statementsRow, args
}
