package object

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Type() Type {
	return NULL_TYPE
}

func (n *Null) Inspect() string {
	return "null"
}

var NULL = NewNull()
