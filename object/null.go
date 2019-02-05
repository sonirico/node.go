package object

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Type() Type {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}
