package goscm

// NOTE: Using reflection this way (converting the type to a string and then
// comparing the string) must be inefficient and not best practise.

import (
	"testing"
	"reflect"
)

func Test_Integer(t *testing.T) {
	// A test with the integer 31337
	in := Make_SCMT(31337)
	if reflect.TypeOf(in).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if in.(*SCMT_Integer).value != 31337 {
		t.Error()
	}
	if in.String() != "31337" {
		t.Error()
	}
}

func Test_String(t *testing.T) {
	// A test with the string "foobar!"
	str := Make_SCMT("foobar!")
	if reflect.TypeOf(str).String() != "*goscm.SCMT_String" {
		t.Error()
	}
	if str.(*SCMT_String).value != "foobar!" {
		t.Error()
	}
	if str.String() != "\"foobar!\"" {
		t.Error()
	}
}

func Test_Pair(t *testing.T) {
	// A test with the singleton ("Foobarrrr!")
	sing := Cons(Make_SCMT("Foobarrrr!"), Make_SCMT(nil))
	if reflect.TypeOf(sing).String() != "*goscm.SCMT_Pair" {
		t.Error()
	}
	if sing.String() != "(\"Foobarrrr!\")" {
		t.Error()
	}
	
	// A test with the pair (4 . 5)
	pair := Cons(Make_SCMT(4), Make_SCMT(5))
	if reflect.TypeOf(pair).String() != "*goscm.SCMT_Pair" {
		t.Error()
	}
	if pair.String() != "(4 . 5)" {
		t.Error()
	}
	
	// A test with the list (1 2 3 4 5 6 7 8 9) aka
	// (1 . (2 . (3 . (4 . (5 . (6 . (7 . (8 . (9 . ())))))))))
	list := Make_SCMT(nil)
	list = Cons(Make_SCMT(9), list)
	list = Cons(Make_SCMT(8), list)
	list = Cons(Make_SCMT(7), list)
	list = Cons(Make_SCMT(6), list)
	list = Cons(Make_SCMT(5), list)
	list = Cons(Make_SCMT(4), list)
	list = Cons(Make_SCMT(3), list)
	list = Cons(Make_SCMT(2), list)
	list = Cons(Make_SCMT(1), list)
	if reflect.TypeOf(list).String() != "*goscm.SCMT_Pair" {
		t.Error()
	}
	if list.String() != "(1 2 3 4 5 6 7 8 9)" {
		t.Error()
	}
}

func Test_Nil(t *testing.T) {
	n := Make_SCMT(nil)
	if reflect.TypeOf(n).String() != "*goscm.SCMT_Nil" {
		t.Error()
	}
	if n.String() != "()" {
		t.Error()
	}
}
