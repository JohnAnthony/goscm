package goscm

type SCMT_Nil struct {}

func (*SCMT_Nil) scm_eval(*SCMT_Environment) SCMT {
	// This is actually an error
	return &SCMT_Nil {}
}

func (*SCMT_Nil) String() string {
	return "()"
}
