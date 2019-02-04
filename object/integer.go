package object

import "strconv"

type Integer struct {
	Value int64
}

func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

func (i *Integer) Type() Type {
	return INT
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}
