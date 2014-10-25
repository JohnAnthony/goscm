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
		t.Error(pair)
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
	if reflect.TypeOf(n).String() != "*goscm.SCMT_Pair" {
		t.Error()
	}
	if n.String() != "()" {
		t.Error()
	}
	if n.(*SCMT_Pair).IsNil() != true {
		t.Error()
	}
}

func Test_Symbol(t *testing.T) {
	s := Make_Symbol("Foo-bar")

	if reflect.TypeOf(s).String() != "*goscm.SCMT_Symbol" {
		t.Error()
	}
	if s.String() != "FOO-BAR" {
		t.Error()
	}
}

func Test_Environment(t *testing.T) {
	env := EnvEmpty(nil)

	if reflect.TypeOf(env).String() != "*goscm.SCMT_Env" {
		t.Error()
	}
	if env.String() != "#<environment>" {
		t.Error()
	}
	
	env.Add(Make_Symbol("derp"), Make_SCMT(9987654))
	ret := env.Find(Make_Symbol("DeRp"))
	if reflect.TypeOf(ret).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if ret.String() != "9987654" {
		t.Error()
	}
}

func Test_Foreign(t *testing.T) {
	f := func (list *SCMT_Pair) SCMT {
		n := Car(list).(*SCMT_Integer).value
		return Make_SCMT(n * n)
	}

	scm_f := Make_Foreign(f)
	if scm_f.String() != "#<foreign function>" {
		t.Error()
	}

	sq := scm_f.Apply(Cons(Make_SCMT(13), Make_SCMT(nil)), EnvEmpty(nil))
	if reflect.TypeOf(sq).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if sq.String() != "169" {
		t.Error()
	}
}

func Test_Foreign_List(t *testing.T) {
	list := SCMT_Nil
	list = Cons(Make_SCMT(11), list)
	list = Cons(Make_SCMT(3), list)
	list = Cons(Make_SCMT(190), list)
	list = Cons(Make_Symbol("+"), list)

	env := EnvEmpty(nil)
	env.BindForeign("+", func (args *SCMT_Pair) SCMT {
		ret := 0
		for ; !args.IsNil(); args = Cdr(args).(*SCMT_Pair) {
			ret += Car(args).(*SCMT_Integer).value
		}
		return Make_SCMT(ret)
	})

	ret := list.scm_eval(env)
	if reflect.TypeOf(ret).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if ret.String() != "204" {
		t.Error()
	}
}

func Test_Special(t *testing.T) {
	env := EnvEmpty(nil)
	env.BindSpecial("quote", func (args *SCMT_Pair) SCMT {
		return args
	})

	list := SCMT_Nil
	list = Cons(Make_Symbol("E"), list)
	list = Cons(Make_Symbol("D"), list)
	list = Cons(Make_Symbol("C"), list)
	list = Cons(Make_Symbol("B"), list)
	list = Cons(Make_Symbol("A"), list)
	list = Cons(Make_Symbol("quote"), list)
	
	ret := list.scm_eval(env)
	if reflect.TypeOf(ret).String() != "*goscm.SCMT_Pair" {
		t.Error()
	}
	if ret.String() != "(A B C D E)" {
		t.Error()
	}
}

// TODO: Test for procedures
