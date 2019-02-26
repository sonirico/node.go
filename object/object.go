package object

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

const (
	INT       Type = "INTEGER"
	BOOL           = "BOOLEAN"
	STRING         = "STRING"
	RETURN         = "RETURN"
	NULL_TYPE      = "NULL"
	ERROR          = "ERROR"
	FUNCTION       = "FUNCTION"
)
