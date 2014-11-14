package goscm

import (
	"testing"
	"reflect"
)

func Test_Symbol(t *testing.T) {
	s := Make_Symbol("Foo-bar")

	if reflect.TypeOf(s) != reflect.TypeOf(&Symbol{}) {
		t.Error()
	}
	if s.String() != "FOO-BAR" {
		t.Error()
	}
}

func Test_Foreign_List(t *testing.T) {
	list := Make_List(
		Make_Symbol("+"),
		Make_SCMT(190),
		Make_SCMT(3),
		Make_SCMT(11),
	)

	env := EnvEmpty(nil)
	env.BindForeign("+", func (args *Pair, env *Environ) (SCMT, error) {
		ret := 0
		for ; !args.IsNil(); args = args.Cdr.(*Pair) {
			ret += args.Car.(*PlainInt).Value
		}
		return Make_SCMT(ret), nil
	})

	ret, err := list.Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ret) != reflect.TypeOf(&PlainInt{}) {
		t.Error()
	}
	if ret.String() != "204" {
		t.Error()
	}
}

func Test_Special(t *testing.T) {
	env := EnvEmpty(nil)
	env.BindSpecial("quote", func (args *Pair, env *Environ) (SCMT, error) {
		return args, nil
	})

	list := Make_List(
		Make_Symbol("quote"),
		Make_Symbol("a"),
		Make_Symbol("b"),
		Make_Symbol("c"),
		Make_Symbol("d"),
		Make_Symbol("e"),
	)
	
	ret, err := list.Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ret) != reflect.TypeOf(&Pair{}) {
		t.Error()
	}
	if ret.String() != "(A B C D E)" {
		t.Error()
	}
}

func Test_Proc(t *testing.T) {
	// We're testing this:
	// ((lambda (n) (* n n)) 123) => 15129
	env := EnvEmpty(nil)
	
	// We have to also provide a multiplication primitive
	scm_multiply := func (args *Pair, env *Environ) (SCMT, error) {
		ret := 1
		for ; !args.IsNil(); args = args.Cdr.(*Pair) {
			ret *= args.Car.(*PlainInt).Value
		}
		return Make_SCMT(ret), nil
	}
	env.BindForeign("*", scm_multiply)

	// args = (n)
	// body = ((* n n))
	args := Make_List(Make_Symbol("n"))
	body := Make_List(Make_List(
		Make_Symbol("*"),
		Make_Symbol("n"),
		Make_Symbol("n"),
	))
	proc, err := Make_Proc(args, body, env)
	if err != nil {
		t.Error(err)
	}
	expr := Make_List(proc, Make_SCMT(123))

	result, err := expr.Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(result) != reflect.TypeOf(&PlainInt{}) {
		t.Error(reflect.TypeOf(result))
	}
	if result.String() != "15129" {
		t.Error(result)
	}
		
}

