package goscm

import (
	"reflect"
)

type SCMT_Pair struct {
	car SCMT
	cdr SCMT
}

func (pair *SCMT_Pair) scm_eval(env *SCMT_Environment) SCMT {
	// COMPLEX!
	return nil
}

func (pair *SCMT_Pair) scm_print() string {
	// NOTE: This reflection to a string followed by a string comparison MUST be
	// a terrible way to do this
	ret := "("
	for {
		ret += pair.car.scm_print()
		if reflect.TypeOf(pair.cdr).String() == "*goscm.SCMT_Nil" {
			ret += ")"
			break
		} else if reflect.TypeOf(pair.cdr).String() != "*goscm.SCMT_Pair" {
			ret += " . "
			ret += pair.cdr.scm_print()
			ret += ")"
			break
		} else { //reflection shows us to have a SCMT_Pair
			ret += " "
			pair = pair.cdr.(*SCMT_Pair)
		}
	}
	return ret
}

// Some standard scheme pair functions provided for use in Go

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
