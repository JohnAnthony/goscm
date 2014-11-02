package goscm

import (
	"reflect"
)

type SCMT_Pair struct {
	Car SCMT
	Cdr SCMT
}

var SCMT_Nil = &SCMT_Pair {}

func (pair *SCMT_Pair) Eval(env *SCMT_Env) SCMT {
	if pair.IsNil() {
		// TODO: Handle this as an error!
		return nil
	}
	proc := pair.Car.Eval(env)
	args := pair.Cdr.(*SCMT_Pair)
	return proc.(SCMT_Func).Apply(args, env)
}

func (pair *SCMT_Pair) String() string {
	if pair.IsNil() { 
		return "()"
	}

	ret := "("
	for {
		ret += pair.Car.String()
		if reflect.TypeOf(pair.Cdr) != reflect.TypeOf(&SCMT_Pair{}) {
			ret = ret + " . " + pair.Cdr.String()
			break
		}
		if pair.Cdr.(*SCMT_Pair).IsNil() {
			break
		} else {
			ret += " "
		}
		pair = pair.Cdr.(*SCMT_Pair)
	}
	return ret + ")"
}

// Some standard scheme pair functions provided for use in Go

func (pair *SCMT_Pair) IsNil() bool {
	return pair == SCMT_Nil
}

func Cons(car SCMT, cdr SCMT) *SCMT_Pair {
	return &SCMT_Pair {
		Car: car,
		Cdr: cdr,
	}
}

func Reverse(pair *SCMT_Pair) *SCMT_Pair {
	ret := SCMT_Nil
	for ; !pair.IsNil(); pair = pair.Cdr.(*SCMT_Pair) {
		ret = Cons(pair.Car, ret)
	}
	return ret
}
