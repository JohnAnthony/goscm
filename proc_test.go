package goscm

import (
	"testing"
	"errors"
)

func Test_Proc(t *testing.T) {
	// We're testing this:
	// ((lambda (n) (* n n)) 123) => 15129
	env := EnvEmpty(nil)
	
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
		return Make_SCMT(ret), nil
	}
	env.Add(Make_Symbol("*"), Make_Foreign(scm_multiply))

	// args = (n)
	// body = ((* n n))
	args := Make_List(Make_Symbol("n"))
	body := Make_List(Make_List(
		Make_Symbol("*"),
		Make_Symbol("n"),
		Make_Symbol("n"),
	))
	proc, err := Make_Proc(args, body, env)
	if err != nil { t.Error(err) }

	expr := Make_List(proc, Make_SCMT(123))
	result, err := expr.Eval(env)
	if err != nil {	t.Error(err) }
	
	if result.String() != "15129" {
		t.Error(result)
	}
}
