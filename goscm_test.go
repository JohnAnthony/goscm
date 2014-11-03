package goscm

import (
	"testing"
	"reflect"
)

func Test_Integer(t *testing.T) {
	// A test with the integer 31337
	in := Make_SCMT(31337)
	if reflect.TypeOf(in) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error()
	}
	if in.(*SCMT_Integer).Value != 31337 {
		t.Error()
	}
	if in.String() != "31337" {
		t.Error()
	}
}

func Test_String(t *testing.T) {
	// A test with the string "foobar!"
	str := Make_SCMT("foobar!")
	if reflect.TypeOf(str) != reflect.TypeOf(&SCMT_String{}) {
		t.Error()
	}
	if str.(*SCMT_String).Value != "foobar!" {
		t.Error()
	}
	if str.String() != "\"foobar!\"" {
		t.Error()
	}
}

func Test_Pair(t *testing.T) {
	// A test with the singleton ("Foobarrrr!")
	sing := Cons(Make_SCMT("Foobarrrr!"), Make_SCMT(nil))
	if reflect.TypeOf(sing) != reflect.TypeOf(&SCMT_Pair{}) {
		t.Error()
	}
	if sing.String() != "(\"Foobarrrr!\")" {
		t.Error()
	}
	
	// A test with the pair (4 . 5)
	pair := Cons(Make_SCMT(4), Make_SCMT(5))
	if reflect.TypeOf(pair) != reflect.TypeOf(&SCMT_Pair{}) {
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
	if reflect.TypeOf(list) != reflect.TypeOf(&SCMT_Pair{}) {
		t.Error()
	}
	if list.String() != "(1 2 3 4 5 6 7 8 9)" {
		t.Error()
	}
	
	// Test Make_List
	list2 := Make_List(
		Make_SCMT(55),
		Make_SCMT(66),
		Make_SCMT(77),
		Make_SCMT(88),
		Make_SCMT(99),
	)
	if reflect.TypeOf(list2) != reflect.TypeOf(&SCMT_Pair{}) {
		t.Error(reflect.TypeOf(list2))
	}
	if list2.String() != "(55 66 77 88 99)" {
		t.Error(list2)
	}
}

func Test_Nil(t *testing.T) {
	n := Make_SCMT(nil)
	if reflect.TypeOf(n) != reflect.TypeOf(&SCMT_Pair{}) {
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

	if reflect.TypeOf(s) != reflect.TypeOf(&SCMT_Symbol{}) {
		t.Error()
	}
	if s.String() != "FOO-BAR" {
		t.Error()
	}
}

func Test_Environment(t *testing.T) {
	env := EnvEmpty(nil)

	if reflect.TypeOf(env) != reflect.TypeOf(&SCMT_Env{}) {
		t.Error()
	}
	if env.String() != "#<environment>" {
		t.Error()
	}
	
	env.Add(Make_Symbol("derp"), Make_SCMT(9987654))
	ret := env.Find(Make_Symbol("DeRp"))
	if reflect.TypeOf(ret) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error()
	}
	if ret.String() != "9987654" {
		t.Error()
	}
}

func Test_Foreign(t *testing.T) {
	f := func (list *SCMT_Pair, env *SCMT_Env) SCMT {
		n := list.Car.(*SCMT_Integer).Value
		return Make_SCMT(n * n)
	}

	scm_f := Make_Foreign(f)
	if scm_f.String() != "#<foreign function>" {
		t.Error()
	}

	sq := scm_f.Apply(Cons(Make_SCMT(13), Make_SCMT(nil)), EnvEmpty(nil))
	if reflect.TypeOf(sq) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error()
	}
	if sq.String() != "169" {
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
	env.BindForeign("+", func (args *SCMT_Pair, env *SCMT_Env) SCMT {
		ret := 0
		for ; !args.IsNil(); args = args.Cdr.(*SCMT_Pair) {
			ret += args.Car.(*SCMT_Integer).Value
		}
		return Make_SCMT(ret)
	})

	ret := list.Eval(env)
	if reflect.TypeOf(ret) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error()
	}
	if ret.String() != "204" {
		t.Error()
	}
}

func Test_Special(t *testing.T) {
	env := EnvEmpty(nil)
	env.BindSpecial("quote", func (args *SCMT_Pair, env *SCMT_Env) SCMT {
		return args
	})

	list := Make_List(
		Make_Symbol("quote"),
		Make_Symbol("a"),
		Make_Symbol("b"),
		Make_Symbol("c"),
		Make_Symbol("d"),
		Make_Symbol("e"),
	)
	
	ret := list.Eval(env)
	if reflect.TypeOf(ret) != reflect.TypeOf(&SCMT_Pair{}) {
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
	scm_multiply := func (args *SCMT_Pair, env *SCMT_Env) SCMT {
		ret := 1
		for ; !args.IsNil(); args = args.Cdr.(*SCMT_Pair) {
			ret *= args.Car.(*SCMT_Integer).Value
		}
		return Make_SCMT(ret)
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
	proc := Make_Proc(args, body, env)
	expr := Make_List(proc, Make_SCMT(123))

	result := expr.Eval(env)
	if reflect.TypeOf(result) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error(reflect.TypeOf(result))
	}
	if result.String() != "15129" {
		t.Error(result)
	}
}
