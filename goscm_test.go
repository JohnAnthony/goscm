package goscm

import (
	"testing"
	"reflect"
)

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

