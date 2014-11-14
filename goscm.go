package goscm

import (
	"reflect"
	"bufio"
	"fmt"
	"strings"
	"errors"
)

type SCMT interface {
	Eval(*Environ) (SCMT, error)
	String() string
}

type SCMT_Func interface {
	Apply(*Pair, *Environ) (SCMT, error)
}

func Make_SCMT(in interface {}) SCMT {
	if in == nil {
		return SCMT_Nil
	}

	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return &PlainInt { Value: in.(int) }
	case reflect.Bool:
		return &Boolean { Value: in.(bool) }
	default:
		// TODO: We probably need to put an error here
		return nil
	}
}

func REPL(in *bufio.Reader, read func(*bufio.Reader) (SCMT, error), env *Environ) {
	for {
		r, err := read(in)
		if err != nil && err.Error() == "EOF" {
			break
		}
		if err != nil {
			fmt.Println("(READ ERROR) " + err.Error())
			continue
		}
		e, err := r.Eval(env)
		if err != nil {
			fmt.Println("(EVAL ERROR) " + err.Error())
			continue
		}
		fmt.Println("=> " + e.String())
	}
}

func EvalStr(str string, read func(*bufio.Reader) (SCMT, error), env *Environ) (SCMT, error) {
	b := bufio.NewReader(strings.NewReader(str))
	r, err := read(b)
	if err != nil && err.Error() != "EOF" {
		return SCMT_Nil, err
	}
	return r.Eval(env)
}

func MapEval(list *Pair, env *Environ) (*Pair, error) {
	if list == SCMT_Nil {
		return SCMT_Nil, nil
	}

	new, err := list.Car.Eval(env)
	if err != nil {
		return SCMT_Nil, nil
	}

	if reflect.TypeOf(list.Cdr) != reflect.TypeOf(&Pair{}) {
		return SCMT_Nil, errors.New("List is not nil-terminated")
	}

	tail, err := MapEval(list.Cdr.(*Pair), env)
	if err != nil {
		return SCMT_Nil, nil
	}

	return Cons(new, tail), nil
}

func IsTrue(s SCMT) bool {
	if reflect.TypeOf(s) == reflect.TypeOf(&Boolean{}) {
		return s.(*Boolean).Value
	}
	return true
}
