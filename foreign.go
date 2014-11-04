package goscm

type SCMT_Foreign struct {
	function func (*SCMT_Pair, *SCMT_Env) (SCMT, error)
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
		val, err := args.Car.Eval(env)
		if err != nil {
			return SCMT_Nil, err
		}
		newargs = Cons(val, newargs)
	}
	newargs = Reverse(newargs)
	return fo.function(newargs, env)
}

func Make_Foreign(f func (*SCMT_Pair, *SCMT_Env) (SCMT, error)) *SCMT_Foreign {
	return &SCMT_Foreign {
		function: f,
	}
}
