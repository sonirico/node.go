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
			strObject, ok := arguments[0].(*String)
			if !ok {
				return NewError(fmt.Sprintf("Type mismatch: Expected %s. Got %s", STRING, arguments[0].Type()))
			}
			return NewInteger(int64(len(strObject.Value)))
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
