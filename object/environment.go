package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Get(ident string) (Object, bool) {
	value, ok := e.store[ident]
	if !ok {
		return NULL, false
	}
	return value, ok
}

func (e *Environment) Set(ident string, value Object) Object {
	e.store[ident] = value
	return value
}
