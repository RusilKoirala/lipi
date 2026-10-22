package evaluator

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/rusilkoirala/lipi/object"
)

// Where bhana prints. The playground (Part 10) will point this somewhere else.
var Out io.Writer = os.Stdout

var builtins = map[string]*object.Builtin{
	// bhana("Namaste", naam)  ->  prints: Namaste Ram
	"bhana": {
		Name: "bhana",
		Fn: func(line int, args ...object.Object) object.Object {
			parts := []string{}
			for _, a := range args {
				parts = append(parts, a.Inspect())
			}
			fmt.Fprintln(Out, strings.Join(parts, " "))
			return NULL
		},
	},

	// lambai("Ram")  ->  3
	"lambai": {
		Name: "lambai",
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, "'lambai' lai 1 ota input chahinchha, %d ota aayo", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(line, "lambai le %s napna sakdaina", args[0].Type())
			}
			return &object.Integer{Value: int64(utf8.RuneCountInString(s.Value))}
		},
	},
}
