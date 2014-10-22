package goscm

type SCMT_String struct {
	value string
}

func (in *SCMT_String) String() string {
	return "\"" + in.value + "\""
}

func (in *SCMT_String) scm_eval(*SCMT_Env) SCMT {
	return in
}
