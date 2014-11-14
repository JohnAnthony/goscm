package goscm

type Foreign struct {
	function func (*Pair, *Environ) (SCMT, error)
}

func (fo *Foreign) Eval(*Environ) (SCMT, error) {
	return fo, nil
}

func (*Foreign) String() string {
	return "#<foreign function>"
}

func (fo *Foreign) Apply(args *Pair, env *Environ) (SCMT, error) {
	newargs := SCMT_Nil
	for ; !args.IsNil(); args = args.Cdr.(*Pair) {
		val, err := args.Car.Eval(env)
		if err != nil {
			return SCMT_Nil, err
		}
		newargs = Cons(val, newargs)
	}
	newargs = Reverse(newargs)
	return fo.function(newargs, env)
}

func Make_Foreign(f func (*Pair, *Environ) (SCMT, error)) *Foreign {
	return &Foreign {
		function: f,
	}
}
