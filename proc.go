package goscm

type ECMT_Proc struct {
	env *SCMT_Env
	body *SCMT_Pair
}

func (p *SCMT_Proc) scm_eval(*SCMT_Env) SCMT {
	return p
}

func (*SCMT_Proc) String() string {
	return "#<procedure>"
}

func (p *SCMT_Proc) Apply(args, *SCMT_Pair, env *SCMT_Env) SCMT {
	return SCMT_Nil
}
