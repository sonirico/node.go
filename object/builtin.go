package object

import "fmt"

type BuiltinFunction func(...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

var builtins = map[string]*Builtin{
	"len": {
		Name: "len",
		Fn: func(arguments ...Object) Object {
			if len(arguments) != 1 {
				return NewError(fmt.Sprintf("Type error: Expected 1 argument. Got %d",
					len(arguments)))
			}
			switch obj := arguments[0].(type) {
			case *String:
				return NewInteger(int64(len(obj.Value)))
			case *Array:
				return NewInteger(int64(len(obj.Items)))
			}
			return NewError(fmt.Sprintf("type mismatch: Expected STRING or ARRAY. Got %s", arguments[0].Type()))
		},
	},
}

func (b *Builtin) Type() Type {
	return BFUNCTION
}
func (b *Builtin) Inspect() string {
	return "__builtin__.len"
}

func LookUpBuiltin(name string) (*Builtin, bool) {
	value, ok := builtins[name]
	return value, ok
}
