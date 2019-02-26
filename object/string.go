package object

import "fmt"

type String struct {
	Value string
}

func NewString(value string) *String {
	return &String{Value: value}
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Inspect() string {
	return fmt.Sprintf("'%s'", s.Value)
}
