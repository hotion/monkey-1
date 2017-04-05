package object

import "fmt"

var Builtins = map[string]Builtin{
	"len": Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("too many arguments. expected=1 got=%d", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Members))}

			}
			return &Error{Message: fmt.Sprintf("unsupported type: %T", args[0])}
		},
	},
	"pop": Builtin{
		Fn: func(args ...Object) Object {
			l := len(args)
			if !(l == 1 || l == 2) {
				return &Error{Message: fmt.Sprintf("too many arguments. expected=1 or 2. got=%d", len(args))}
			}
			switch obj := args[l-1].(type) {
			case *Array:
				if l == 1 {
					popped, shifted := obj.Members[0], obj.Members[1:]
					obj.Members = shifted
					return popped
				} else {
					idx := args[0].(*Integer).Value
					popped := obj.Members[idx]
					obj.Members = append(obj.Members[:idx], obj.Members[idx+1:]...)
					return popped
				}
			default:
				return &Error{Message: fmt.Sprintf("unsupportedtype %T", args[0])}
			}
		},
	},
	"push": Builtin{
		Fn: func(args ...Object) Object {
			l := len(args)
			if l != 2 {
				return &Error{Message: fmt.Sprintf("too many arguments. expected=1. got=%d", len(args))}
			}
			switch obj := args[l-1].(type) {
			case *Array:
				obj.Members = append(obj.Members, args[0])
				return obj
			default:
				return &Error{Message: fmt.Sprintf("unsupportedtype %T", args[0])}
			}
		},
	},
}
