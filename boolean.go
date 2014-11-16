package goscm

type Boolean struct {
	Value bool 
}

func (b *Boolean) String() string {
	if b.Value {
		return "#t"
	}
	return "#f" 
}

func (in *Boolean) Eval(*Environ) (SCMT, error) {
	return in, nil
}

func NewBoolean(b bool) *Boolean {
	return &Boolean { Value: b }
}
