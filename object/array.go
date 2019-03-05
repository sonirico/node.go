package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Items []Object
}

func NewArray(items []Object) *Array {
	return &Array{Items: items}
}

func (a *Array) Type() Type {
	return ARRAY
}
func (a *Array) Inspect() string {
	var out bytes.Buffer
	var items []string

	for _, item := range a.Items {
		items = append(items, item.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")

	return out.String()
}
