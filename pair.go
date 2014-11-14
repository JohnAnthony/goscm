package goscm

import (
	"reflect"
	"errors"
)

type Pair struct {
	Car SCMT
	Cdr SCMT
}

var SCMT_Nil = &Pair {}

func (pair *Pair) Eval(env *Environ) (SCMT, error) {
	if pair.IsNil() {
		return SCMT_Nil, errors.New("Cannot eval empty list")
	}
	proc, err := pair.Car.Eval(env)
	if err != nil {
		return SCMT_Nil, err
	}
	args := pair.Cdr.(*Pair)
	return proc.(Func).Apply(args, env)
}

func (pair *Pair) String() string {
	if pair.IsNil() { 
		return "()"
	}

	ret := "("
	for {
		ret += pair.Car.String()
		if reflect.TypeOf(pair.Cdr) != reflect.TypeOf(&Pair{}) {
			ret = ret + " . " + pair.Cdr.String()
			break
		}
		if pair.Cdr.(*Pair).IsNil() {
			break
		} else {
			ret += " "
		}
		pair = pair.Cdr.(*Pair)
	}
	return ret + ")"
}

func Cast_Pair(scm SCMT) (*Pair, error) {
	if reflect.TypeOf(scm) != reflect.TypeOf(&Pair{}) {
		return SCMT_Nil, errors.New("Cast failed: Pair")
	}
	return scm.(*Pair), nil
}

// Helpers

func Make_List(args ...SCMT) *Pair {
	list := SCMT_Nil
	for i := len(args) - 1; i >= 0; i-- {
		list = Cons(args[i], list)
	}
	return list
}

func (p *Pair) ToSlice() ([]SCMT, error) {
	var s []SCMT
	var err error

	for p != SCMT_Nil {
		s = append(s, p.Car)
		p, err = Cast_Pair(p.Cdr)
		if err != nil {
			return s, errors.New("Non-nil terminated list")
		}
	}
	return s, nil
}

func (pair *Pair) IsNil() bool {
	return pair == SCMT_Nil
}

func Cons(car SCMT, cdr SCMT) *Pair {
	return &Pair {
		Car: car,
		Cdr: cdr,
	}
}

func Reverse(pair *Pair) *Pair {
	ret := SCMT_Nil
	for ; !pair.IsNil(); pair = pair.Cdr.(*Pair) {
		ret = Cons(pair.Car, ret)
	}
	return ret
}
