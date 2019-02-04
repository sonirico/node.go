package object

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

var (
	INT Type = "INTEGER"
)
