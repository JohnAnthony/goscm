package goscm

import (
	"testing"
)

func Test_Environment(t *testing.T) {
	// Making a top-level environment
	top := EnvEmpty(nil)
	if top.String() != "#<environment>" {
		t.Error()
	}

	// Adding to an environment and retrieving
	top.Add(Make_Symbol("foo"), Make_PlainInt(9987654))
	ret, err := top.Find(Make_Symbol("BAr"))
	if err != nil {	t.Error(err) }
	ans, ok := ret.(*PlainInt)
	if !ok { t.Error(ok) }
	if ans.String() != "9987654" {
		t.Error()
	}
	
	// Making a child environment and retrieving from it
	env := EnvEmpty(top)
	env.Add(Make_Symbol("bar"), Make_PlainInt(9987666))
	ret, err = env.Find(Make_Symbol("FOO"))
	if err != nil { t.Error(err) }
	ans, ok = ret.(*PlainInt)
	if !ok { t.Error(ok) }
	if ans.String() != "9987666" {
		t.Error()
	}
	
	// Retrieving from a parent environment from a child
	ret, err = env.Find(Make_Symbol("BAr"))
	if err != nil { t.Error(err) }
	ans, ok = ret.(*PlainInt)
	if !ok { t.Error(ok) }
	if ans.String() != "9987654" {
		t.Error()
	}
}
