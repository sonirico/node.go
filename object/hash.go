package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (s *String) HashKey() HashKey {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: hash.Sum64()}
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
	var pairsString []string

	for _, pair := range h.Pairs {
		pairsString = append(pairsString, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairsString, ", "))
	out.WriteString("}")

	return out.String()
}
