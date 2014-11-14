package goscm

import "errors"

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
	var ok bool
	var ret *Pair

	for ret = SCMT_Nil; args != SCMT_Nil; args, ok = args.Cdr.(*Pair) {
		if !ok { // This is a dotted list
			return SCMT_Nil, errors.New("Got a dotted list. How to handle?")
		}

		val, err := args.Car.Eval(env)
		if err != nil {	return SCMT_Nil, err }

		ret = Cons(val, ret)
	}

	ret = Reverse(ret)
	return fo.function(ret, env)
}

func Make_Foreign(f func (*Pair, *Environ) (SCMT, error)) *Foreign {
	return &Foreign {
		function: f,
	}
}
