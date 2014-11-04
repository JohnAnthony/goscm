package goscm

import (
	"reflect"
	"bufio"
	"fmt"
)

type SCMT interface {
	Eval(*SCMT_Env) (SCMT, error)
	String() string
}

type SCMT_Func interface {
	Apply(*SCMT_Pair, *SCMT_Env) (SCMT, error)
}

func Make_SCMT(in interface {}) SCMT {
	if in == nil {
		return SCMT_Nil
	}

	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return &SCMT_Integer { Value: in.(int) }
	case reflect.String:
		return &SCMT_String { Value: in.(string) }
	default:
		// TODO: We probably need to put an error here
		return nil
	}
}

func REPL(in *bufio.Reader, read func(*bufio.Reader) (SCMT, error), env *SCMT_Env) {
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
