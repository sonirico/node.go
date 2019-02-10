package object

type Return struct {
	Value Object
}

func (r *Return) Type() Type {
	return RETURN
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}
