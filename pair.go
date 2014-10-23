package goscm

import (
	"reflect"
)

type SCMT_Pair struct {
	car SCMT
	cdr SCMT
}

var SCMT_Nil = &SCMT_Pair {}

func (pair *SCMT_Pair) scm_eval(env *SCMT_Env) SCMT {
	if pair.IsNil() {
		// TODO: Handle this as an error!
		return nil
	}

	proc := Car(pair).(SCMT_Func)
	args := SCMT_Nil
	for pair = Cdr(pair).(*SCMT_Pair); !pair.IsNil(); pair = Cdr(pair).(*SCMT_Pair) {
		args = Cons(pair.car.scm_eval(env), args)
		if reflect.TypeOf(Cdr(pair)).String() != "*goscm.SCMT_Pair" {
			args = Cons(pair.cdr.scm_eval(env), args)
			break
		}
	}
	return proc.Apply(args)
}

func (pair *SCMT_Pair) String() string {
	if pair.IsNil() { 
		return "()"
	}

	ret := "("
	for {
		ret += pair.car.String()
		if reflect.TypeOf(pair.cdr).String() != "*goscm.SCMT_Pair" {
			ret = ret + " . " + pair.cdr.String()
			break
		}
		if pair.cdr.(*SCMT_Pair).IsNil() {
			break
		} else {
			ret += " "
		}
		pair = pair.cdr.(*SCMT_Pair)
	}
	return ret + ")"
}

// Some standard scheme pair functions provided for use in Go

func (pair *SCMT_Pair) IsNil() bool {
	return pair == SCMT_Nil
}

func Cons(car SCMT, cdr SCMT) *SCMT_Pair {
	return &SCMT_Pair {
		car: car,
		cdr: cdr,
	}
}

func Car(pair *SCMT_Pair) SCMT {
	return pair.car
}

func Cdr(pair *SCMT_Pair) SCMT {
	return pair.cdr
}
