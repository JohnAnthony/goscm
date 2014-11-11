package goscm

type SCMT_Special struct {
	function func (*SCMT_Pair, *SCMT_Env) (SCMT, error)
	Expand func (*SCMT_Pair, *SCMT_Env) (SCMT, error) // Used for TCO
}

func (sp *SCMT_Special) Eval(*SCMT_Env) (SCMT, error) {
	return sp, nil
}

func (*SCMT_Special) String() string {
	return "#<special form>"
}

func (sp *SCMT_Special) Apply(args *SCMT_Pair, env *SCMT_Env) (SCMT, error) {
	return sp.function(args, env)
}

func Make_Special(f func (*SCMT_Pair, *SCMT_Env) (SCMT, error)) *SCMT_Special {
	return &SCMT_Special {
		function: f,
		Expand: nil,
	}
}

func Make_SpecialExpandable(f func (*SCMT_Pair, *SCMT_Env) (SCMT, error), ex func (*SCMT_Pair, *SCMT_Env) (SCMT, error)) *SCMT_Special {
	return &SCMT_Special {
		function: f,
		Expand: ex,
	}
}
