package goscm

import (
	"reflect"
	"errors"
)

type SCMT_Pair struct {
	Car SCMT
	Cdr SCMT
}

var SCMT_Nil = &SCMT_Pair {}

func (pair *SCMT_Pair) Eval(env *SCMT_Env) (SCMT, error) {
	if pair.IsNil() {
		return SCMT_Nil, errors.New("Cannot eval empty list")
	}
	proc, err := pair.Car.Eval(env)
	if err != nil {
		return SCMT_Nil, err
	}
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

func Cast_Pair(scm SCMT) (*SCMT_Pair, error) {
	if reflect.TypeOf(scm) != reflect.TypeOf(&SCMT_Pair{}) {
		return SCMT_Nil, errors.New("Cast failed: Pair")
	}
	return scm.(*SCMT_Pair), nil
}

// Helpers

func Make_List(args ...SCMT) *SCMT_Pair {
	list := SCMT_Nil
	for i := len(args) - 1; i >= 0; i-- {
		list = Cons(args[i], list)
	}
	return list
}

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
