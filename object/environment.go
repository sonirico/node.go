package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(ident string) (Object, bool) {
	value, ok := e.store[ident]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(ident)
	}
	return value, ok
}

func (e *Environment) Set(ident string, value Object) Object {
	e.store[ident] = value
	return value
}
