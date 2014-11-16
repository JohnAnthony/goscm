package goscm

import (
	"errors"
	"reflect"
)

type Pair struct {
	Car SCMT
	Cdr SCMT
}

var SCM_Nil = &Pair {}

func (pair *Pair) Eval(env *Environ) (SCMT, error) {
	if pair == SCM_Nil {
		return SCM_Nil, errors.New("Cannot eval empty list")
	}

	proc, err := pair.Car.Eval(env)
	if err != nil {	return SCM_Nil, err	}

	args, ok := pair.Cdr.(*Pair)
	if !ok {
		return SCM_Nil, errors.New("Cannot eval tupple")
	}
	
	f, ok := proc.(Func)
	if !ok {
		return SCM_Nil, errors.New("Attempting to apply non-function")
	}

	return f.Apply(args, env)
}

func (pair *Pair) String() string {
	ret := "("
	for pair != SCM_Nil {
		ret += pair.Car.String()

		if pair.Cdr == SCM_Nil {
			break
		}

		if reflect.TypeOf(pair.Cdr) != reflect.TypeOf(&Pair{}) {
			ret += " . "
			ret += pair.Cdr.String()
			break
		}

		pair = pair.Cdr.(*Pair)
		ret += " "
	}
	return ret + ")"
}

func NewList(args ...SCMT) *Pair {
	list := SCM_Nil
	for i := len(args) - 1; i >= 0; i-- {
		list = Cons(args[i], list)
	}
	return list
}

// Helpers

func (p *Pair) ToSlice() ([]SCMT, error) {
	var s []SCMT
	var ok bool

	for p != SCM_Nil {
		s = append(s, p.Car)
		p, ok = p.Cdr.(*Pair)
		if !ok {
			s = append(s, p.Cdr)
			break
		}
	}
	return s, nil
}

func (list *Pair) MapEval(env *Environ) (*Pair, error) {
	ret := SCM_Nil
	listt, err := list.ToSlice()
	if err != nil { return SCM_Nil, err }
	
	for i := len(listt) - 1; i >= 0; i-- {
		newval, err := listt[i].Eval(env)
		if err != nil { return SCM_Nil, err }

		ret = Cons(newval, ret)
	}

	return ret, nil
}

func Cons(car SCMT, cdr SCMT) *Pair {
	return &Pair {
		Car: car,
		Cdr: cdr,
	}
}

func Reverse(pair *Pair) (*Pair, error) {
	var ok bool
	ret := SCM_Nil

	for pair != SCM_Nil {
		ret = Cons(pair.Car, ret)
		pair, ok = pair.Cdr.(*Pair) 
		if !ok {
			return SCM_Nil, errors.New("Attempt to reverse dotted list")
		}
	}
	return ret, nil
}
