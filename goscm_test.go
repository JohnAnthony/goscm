package goscm

import (
	"testing"
	"reflect"
)

func Test_Symbol(t *testing.T) {
	s := NewSymbol("Foo-bar")

	if reflect.TypeOf(s) != reflect.TypeOf(&Symbol{}) {
		t.Error()
	}
	if s.String() != "FOO-BAR" {
		t.Error()
	}
}

func Test_Foreign_List(t *testing.T) {
	list := NewList(
		NewSymbol("+"),
		NewSCMT(190),
		NewSCMT(3),
		NewSCMT(11),
	)

	env := EnvEmpty(nil)
	env.BindForeign("+", func (args *Pair, env *Environ) (SCMT, error) {
		ret := 0
		for ; !args.IsNil(); args = args.Cdr.(*Pair) {
			ret += args.Car.(*PlainInt).Value
		}
		return NewSCMT(ret), nil
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

	list := NewList(
		NewSymbol("quote"),
		NewSymbol("a"),
		NewSymbol("b"),
		NewSymbol("c"),
		NewSymbol("d"),
		NewSymbol("e"),
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
