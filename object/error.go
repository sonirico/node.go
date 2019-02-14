package object

type Error struct {
	Message string
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
