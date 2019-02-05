package object

type Boolean struct {
	Value bool
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{Value: value}
}

func (b *Boolean) Type() Type {
	return BOOL
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}
