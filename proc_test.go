package goscm

import (
	"testing"
	"errors"
)

func Test_Proc(t *testing.T) {
	// We're testing this:
	// ((lambda (n) (* n n)) 123) => 15129
	env := NewEnv(nil)
	
	// We have to also provide a multiplication primitive
	scm_multiply := func (args *Pair, env *Environ) (SCMT, error) {
		ret := 1
		for args != SCM_Nil {
			pl, ok := args.Car.(*PlainInt)
			if !ok {
				return SCM_Nil, errors.New("Expected PlainInt argument")
			}

			ret *= pl.Value

			args, ok = args.Cdr.(*Pair) 
			if !ok {
				return SCM_Nil, errors.New("Can't operate on dotted list")
			}
		}
		return NewSCMT(ret), nil
	}
	env.Add(NewSymbol("*"), NewForeign(scm_multiply))

	// args = (n)
	// body = ((* n n))
	args := NewList(NewSymbol("n"))
	body := NewList(NewList(
		NewSymbol("*"),
		NewSymbol("n"),
		NewSymbol("n"),
	))
	proc, err := NewProc(args, body, env)
	if err != nil { t.Error(err) }

	expr := NewList(proc, NewSCMT(123))
	result, err := expr.Eval(env)
	if err != nil {	t.Error(err) }
	
	if result.String() != "15129" {
		t.Error(result)
	}
}
