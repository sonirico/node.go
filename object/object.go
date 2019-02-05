package object

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

const (
	INT  Type = "INTEGER"
	BOOL      = "BOOLEAN"
	NULL      = "NULL"
)
