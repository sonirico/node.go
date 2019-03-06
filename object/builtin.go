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
		Fn:   Len,
	},
	"head": {
		Name: "head",
		Fn:   Head,
	},
	"foot": {
		Name: "foot",
		Fn:   Foot,
	},
	"tail": {
		Name: "tail",
		Fn:   Tail,
	},
}

func (b *Builtin) Type() Type {
	return BFUNCTION
}
func (b *Builtin) Inspect() string {
	return "__builtin__." + b.Name
}

func LookUpBuiltin(name string) (*Builtin, bool) {
	value, ok := builtins[name]
	return value, ok
}

// LEN

func Len(arguments ...Object) Object {
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
}

// HEAD

func Head(arguments ...Object) Object {
	if len(arguments) != 1 {
		return NewError(fmt.Sprintf("type error: Expected 1 argument. Got %d",
			len(arguments)))
	}
	switch obj := arguments[0].(type) {
	case *Array:
		if len(obj.Items) < 1 {
			return NULL
		}
		return obj.Items[0]
	}
	return NewError(fmt.Sprintf("type mismatch: Expected ARRAY. Got %s", arguments[0].Type()))
}

// FOOT

func Foot(arguments ...Object) Object {
	if len(arguments) != 1 {
		return NewError(fmt.Sprintf("type error: Expected 1 argument. Got %d",
			len(arguments)))
	}
	switch obj := arguments[0].(type) {
	case *Array:
		length := len(obj.Items)
		if length < 1 {
			return NULL
		}
		return obj.Items[length-1]
	}
	return NewError(fmt.Sprintf("type mismatch: Expected ARRAY. Got %s", arguments[0].Type()))
}

// TAIL

func Tail(arguments ...Object) Object {
	if len(arguments) != 1 {
		return NewError(fmt.Sprintf("type error: Expected 1 argument. Got %d",
			len(arguments)))
	}
	switch obj := arguments[0].(type) {
	case *Array:
		length := len(obj.Items)
		if length < 1 {
			return NULL
		}
		items := make([]Object, length-1, length-1)
		copy(items, obj.Items[1:length])
		return NewArray(items)
	}
	return NewError(fmt.Sprintf("type mismatch: Expected ARRAY. Got %s", arguments[0].Type()))
}
