package goscm

type SCMT_Foreign struct {
	function func (*SCMT_Pair, *SCMT_Env) SCMT
}

func (fo *SCMT_Foreign) Eval(*SCMT_Env) (SCMT, error) {
	return fo, nil
}

func (*SCMT_Foreign) String() string {
	return "#<foreign function>"
}

func (fo *SCMT_Foreign) Apply(args *SCMT_Pair, env *SCMT_Env) (SCMT, error) {
	newargs := SCMT_Nil
	for ; !args.IsNil(); args = args.Cdr.(*SCMT_Pair) {
		val, _ := args.Car.Eval(env)
		newargs = Cons(val, newargs)
	}
	newargs = Reverse(newargs)
	return fo.function(newargs, env), nil
}

func Make_Foreign(f func (*SCMT_Pair, *SCMT_Env) SCMT) *SCMT_Foreign {
	return &SCMT_Foreign {
		function: f,
	}
}
