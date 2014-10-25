package goscm

type SCMT_Special struct {
	function func (*SCMT_Pair, *SCMT_Env) SCMT
}

func (sp *SCMT_Special) scm_eval(*SCMT_Env) SCMT {
	return sp
}

func (*SCMT_Special) String() string {
	return "#<special form>"
}

func (sp *SCMT_Special) Apply(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return sp.function(args, env) // TODO
}

func Make_Special(f func (*SCMT_Pair, *SCMT_Env) SCMT) *SCMT_Special {
	return &SCMT_Special {
		function: f,
	}
}
