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
	var err error
	var ret *Pair

	for ret = SCM_Nil; args != SCM_Nil; args, ok = args.Cdr.(*Pair) {
		if !ok { // This is a dotted list
			return SCM_Nil, errors.New("Got a dotted list. How to handle?")
		}

		val, err := args.Car.Eval(env)
		if err != nil {	return SCM_Nil, err }

		ret = Cons(val, ret)
	}

	ret, err = Reverse(ret)
	if err != nil { return SCM_Nil, err }

	return fo.function(ret, env)
}

func Make_Foreign(f func (*Pair, *Environ) (SCMT, error)) *Foreign {
	return &Foreign {
		function: f,
	}
}
