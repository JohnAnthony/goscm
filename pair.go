package goscm

import "reflect"

type SCMT_Pair struct {
	car SCMT
	cdr SCMT
}

func (pair *SCMT_Pair) scm_eval(env *SCMT_Environment) SCMT {
	// COMPLEX!
	return nil
}

func (pair *SCMT_Pair) scm_print() string {
	ret := "("
	for {
		ret += pair.car.scm_print()
		if pair.cdr == nil {
			ret += ")"
			break
		} else if reflect.TypeOf(pair.cdr).Name() != "SCMT_Pair" {
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
