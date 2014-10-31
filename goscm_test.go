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
	f := func (list *SCMT_Pair, env *SCMT_Env) SCMT {
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
	env.BindForeign("+", func (args *SCMT_Pair, env *SCMT_Env) SCMT {
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
	env.BindSpecial("quote", func (args *SCMT_Pair, env *SCMT_Env) SCMT {
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

func Test_EnvSimple(t *testing.T) {
	env := EnvSimple()

	define_expr := SCMT_Nil
	define_expr = Cons(Make_SCMT(1234), define_expr)
	define_expr = Cons(Make_Symbol("a"), define_expr)
	define_expr = Cons(Make_Symbol("define"), define_expr)
	define_expr.scm_eval(env)
	a_result := Make_Symbol("a").scm_eval(env)
	if reflect.TypeOf(a_result).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if a_result.String() != "1234" {
		t.Error()
	}
	
	add_expr := SCMT_Nil
	add_expr = Cons(Make_SCMT(1000), add_expr)
	add_expr = Cons(Make_SCMT(222), add_expr)
	add_expr = Cons(Make_Symbol("a"), add_expr)
	add_expr = Cons(Make_SCMT(321), add_expr)
	add_expr = Cons(Make_Symbol("+"), add_expr)
	add_result := add_expr.scm_eval(env)
	if reflect.TypeOf(add_result).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if add_result.String() != "2777" {
		t.Error()
	}
	
	sub_expr := SCMT_Nil
	sub_expr = Cons(Make_SCMT(90), sub_expr)
	sub_expr = Cons(Make_SCMT(9), sub_expr)
	sub_expr = Cons(Make_SCMT(100), sub_expr)
	sub_expr = Cons(Make_Symbol("-"), sub_expr)
	sub_result := sub_expr.scm_eval(env)
	if reflect.TypeOf(sub_result).String() != "*goscm.SCMT_Integer" {
		t.Error()
	}
	if sub_result.String() != "1" {
		t.Error()
	}

	mult_expr := SCMT_Nil
	mult_expr = Cons(Make_SCMT(5), mult_expr)
	mult_expr = Cons(Make_SCMT(4), mult_expr)
	mult_expr = Cons(Make_SCMT(3), mult_expr)
	mult_expr = Cons(Make_SCMT(2), mult_expr)
	mult_expr = Cons(Make_SCMT(1), mult_expr)
	mult_expr = Cons(Make_Symbol("*"), mult_expr)
	mult_result := mult_expr.scm_eval(env)
	if reflect.TypeOf(mult_result).String() != "*goscm.SCMT_Integer" {
		t.Error(reflect.TypeOf(mult_result))
	}
	if mult_result.String() != "120" {
		t.Error(mult_result)
	}
	
	div_expr := SCMT_Nil
	div_expr = Cons(Make_SCMT(2), div_expr)
	div_expr = Cons(Make_SCMT(3), div_expr)
	div_expr = Cons(Make_SCMT(66), div_expr)
	div_expr = Cons(Make_Symbol("/"), div_expr)
	div_result := div_expr.scm_eval(env)
	if reflect.TypeOf(div_result).String() != "*goscm.SCMT_Integer" {
		t.Error(reflect.TypeOf(div_result))
	}
	if div_result.String() != "11" {
		t.Error(div_result)
	}
	
	// Building this list: (car (cons "FOO" "BAR"))
	car_expr := SCMT_Nil
	car_expr = Cons(Make_SCMT("BAR"), car_expr)
	car_expr = Cons(Make_SCMT("FOO"), car_expr)
	car_expr = Cons(Make_Symbol("cons"), car_expr)
	car_expr = Cons(car_expr, SCMT_Nil)
	car_expr = Cons(Make_Symbol("car"), car_expr)
	car_result := car_expr.scm_eval(env)
	if reflect.TypeOf(car_result).String() != "*goscm.SCMT_String" {
		t.Error(reflect.TypeOf(car_result))
	}
	if car_result.String() != "\"FOO\"" {
		t.Error(car_result)
	}

	// Building this list: (cdr (cons "FOO" "BAR"))
	cdr_expr := SCMT_Nil
	cdr_expr = Cons(Make_SCMT("BAR"), cdr_expr)
	cdr_expr = Cons(Make_SCMT("FOO"), cdr_expr)
	cdr_expr = Cons(Make_Symbol("cons"), cdr_expr)
	cdr_expr = Cons(cdr_expr, SCMT_Nil)
	cdr_expr = Cons(Make_Symbol("cdr"), cdr_expr)
	cdr_result := cdr_expr.scm_eval(env)
	if reflect.TypeOf(cdr_result).String() != "*goscm.SCMT_String" {
		t.Error(reflect.TypeOf(cdr_result))
	}
	if cdr_result.String() != "\"BAR\"" {
		t.Error(cdr_result)
	}
	
	cons_expr := SCMT_Nil
	cons_expr = Cons(Make_SCMT(5), cons_expr)
	cons_expr = Cons(Make_SCMT(2), cons_expr)
	cons_expr = Cons(Make_Symbol("cons"), cons_expr)
	cons_result := cons_expr.scm_eval(env)
	if reflect.TypeOf(cons_result).String() != "*goscm.SCMT_Pair" {
		t.Error(reflect.TypeOf(cons_result))
	}
	if cons_result.String() != "(2 . 5)" {
		t.Error(cons_result)
	}
	
	// This is a little complex - we're actually making a list whose first
	// element is the symbol "quote", whose second element is a complicated list
	// and which is then terminated by SCMT_Nil.
	quote_expr := SCMT_Nil
	quote_expr = Cons(Make_SCMT(9), quote_expr)
	quote_expr = Cons(Make_SCMT(8), quote_expr)
	quote_expr = Cons(Make_SCMT(7), quote_expr)
	quote_expr = Cons(Make_Symbol("honk"), quote_expr)
	quote_expr = Cons(Make_SCMT(6), quote_expr)
	quote_expr = Cons(Make_SCMT(5), quote_expr)
	quote_expr = Cons(Make_SCMT(4), quote_expr)
	quote_expr = Cons(Make_SCMT(3), quote_expr)
	quote_expr = Cons(Make_Symbol("squeak"), quote_expr)
	quote_expr = Cons(Make_SCMT(2), quote_expr)
	quote_expr = Cons(Make_SCMT(1), quote_expr)
	quote_expr = Cons(Make_Symbol("roar"), quote_expr)
	quote_expr = Cons(quote_expr, SCMT_Nil)
	quote_expr = Cons(Make_Symbol("quote"), quote_expr)
	quote_result := quote_expr.scm_eval(env)
	if reflect.TypeOf(quote_result).String() != "*goscm.SCMT_Pair" {
		t.Error(reflect.TypeOf(quote_result))
	}
	if quote_result.String() != "(ROAR 1 2 SQUEAK 3 4 5 6 HONK 7 8 9)" {
		t.Error(quote_result)
	}

	// let
	// begin
	// lambda
}

// TODO: Test for procedures
