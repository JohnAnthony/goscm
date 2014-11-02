package goscm

type SCMT_String struct {
	Value string
}

func (in *SCMT_String) String() string {
	return "\"" + in.Value + "\""
}

func (in *SCMT_String) Eval(*SCMT_Env) SCMT {
	return in
}
