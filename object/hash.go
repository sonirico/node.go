package object

import (
	"bytes"
	"hash/fnv"
)

type Hashable interface {
	HashKey() *HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

func (i *Integer) HashKey() *HashKey {
	return &HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (b *Boolean) HashKey() *HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return &HashKey{Type: b.Type(), Value: value}
}

func (s *String) HashKey() *HashKey {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(s.Value))
	return &HashKey{Type: s.Type(), Value: hash.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() Type {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")

	for _, pair := range h.Pairs {
		out.WriteString(pair.Key.Inspect())
		out.WriteString(": ")
		out.WriteString(pair.Value.Inspect())
	}

	out.WriteString("}")

	return out.String()
}
