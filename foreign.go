package goscm

type SCMT_Foreign struct {
	function func (*SCMT_Pair, *SCMT_Env) SCMT
}

func (fo *SCMT_Foreign) scm_eval(*SCMT_Env) SCMT {
	return fo
}

func (*SCMT_Foreign) String() string {
	return "#<foreign function>"
}

func (fo *SCMT_Foreign) Apply(args *SCMT_Pair, env *SCMT_Env) SCMT {
	newargs := SCMT_Nil
	for ; !args.IsNil(); args = Cdr(args).(*SCMT_Pair) {
		newargs = Cons(Car(args).scm_eval(env), newargs)
	}
	newargs = Reverse(newargs)
	return fo.function(newargs, env)
}

func Make_Foreign(f func (*SCMT_Pair, *SCMT_Env) SCMT) *SCMT_Foreign {
	return &SCMT_Foreign {
		function: f,
	}
}
