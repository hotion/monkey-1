package eval

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

var builtins map[string]*Builtin

func init() {
	builtins = map[string]*Builtin{
		"int": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch input := args[0].(type) {
				case *Integer:
					return input
				case *String:
					n, err := strconv.Atoi(input.Value)
					if err != nil {
						return newError(INPUTERROR, "STRING: "+input.Value, "int")
					}
					return &Integer{Value: int64(n)}
				}
				return newError(INPUTERROR, args[0].Type(), "int")
			},
		},
		"str": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch input := args[0].(type) {
				case *String:
					return input
				default:
					return &String{Value: input.Inspect()}
				}
				return newError(INPUTERROR, args[0].Type(), "str")
			},
		},
		"len": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Members))}

				}
				return newError(NOMETHODERROR, "len", args[0].Type())
			},
		},
		"methods": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				methods := &Array{}
				t := reflect.TypeOf(args[0])
				for i := 0; i < t.NumMethod(); i++ {
					m := t.Method(i).Name
					if !(m == "Type" || m == "CallMethod" || m == "HashKey" || m == "Inspect") {
						methods.Members = append(methods.Members, &String{Value: strings.ToLower(m)})
					}
				}
				return methods
			},
		},
		"puts": &Builtin{
			Fn: func(args ...Object) Object {
				fmt.Println(args[0].Inspect())
				return NULL
			},
		},
		"type": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				return &String{Value: fmt.Sprintf("%s", args[0].Type())}
			},
		},
	}
}
