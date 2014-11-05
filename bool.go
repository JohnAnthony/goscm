package goscm

type SCMT_Bool struct {
	Value bool 
}

func (b *SCMT_Bool) String() string {
	if b.Value {
		return "#t"
	}
	return "#f" 
}

func (in *SCMT_Bool) Eval(*SCMT_Env) (SCMT, error) {
	return in, nil
}
