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
	if reflect.TypeOf(str) != reflect.TypeOf(&SCMT_String{}) {
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
		n := Car(list).(*SCMT_Integer).value
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

	list := SCMT_Nil
	list = Cons(Make_Symbol("E"), list)
	list = Cons(Make_Symbol("D"), list)
	list = Cons(Make_Symbol("C"), list)
	list = Cons(Make_Symbol("B"), list)
	list = Cons(Make_Symbol("A"), list)
	list = Cons(Make_Symbol("quote"), list)
	
	ret := list.Eval(env)
	if reflect.TypeOf(ret) != reflect.TypeOf(&SCMT_Pair{}) {
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
	define_expr.Eval(env)
	a_result := Make_Symbol("a").Eval(env)
	if reflect.TypeOf(a_result) != reflect.TypeOf(&SCMT_Integer{}) {
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
	add_result := add_expr.Eval(env)
	if reflect.TypeOf(add_result) != reflect.TypeOf(&SCMT_Integer{}) {
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
	sub_result := sub_expr.Eval(env)
	if reflect.TypeOf(sub_result) != reflect.TypeOf(&SCMT_Integer{}) {
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
	mult_result := mult_expr.Eval(env)
	if reflect.TypeOf(mult_result) != reflect.TypeOf(&SCMT_Integer{}) {
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
	div_result := div_expr.Eval(env)
	if reflect.TypeOf(div_result) != reflect.TypeOf(&SCMT_Integer{}) {
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
	car_result := car_expr.Eval(env)
	if reflect.TypeOf(car_result) != reflect.TypeOf(&SCMT_String{}) {
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
	cdr_result := cdr_expr.Eval(env)
	if reflect.TypeOf(cdr_result) != reflect.TypeOf(&SCMT_String{}) {
		t.Error(reflect.TypeOf(cdr_result))
	}
	if cdr_result.String() != "\"BAR\"" {
		t.Error(cdr_result)
	}
	
	cons_expr := SCMT_Nil
	cons_expr = Cons(Make_SCMT(5), cons_expr)
	cons_expr = Cons(Make_SCMT(2), cons_expr)
	cons_expr = Cons(Make_Symbol("cons"), cons_expr)
	cons_result := cons_expr.Eval(env)
	if reflect.TypeOf(cons_result) != reflect.TypeOf(&SCMT_Pair{}) {
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
	quote_result := quote_expr.Eval(env)
	if reflect.TypeOf(quote_result) != reflect.TypeOf(&SCMT_Pair{}) {
		t.Error(reflect.TypeOf(quote_result))
	}
	if quote_result.String() != "(ROAR 1 2 SQUEAK 3 4 5 6 HONK 7 8 9)" {
		t.Error(quote_result)
	}

	// This is a little complex as well. We need to make a list that looks like
	// (begin
	//  (define x 10)
	//  (* x 3))
	// And check that it returns 30
	begin_subexpr1 := SCMT_Nil
	begin_subexpr1 = Cons(Make_SCMT(10), begin_subexpr1)
	begin_subexpr1 = Cons(Make_Symbol("x"), begin_subexpr1)
	begin_subexpr1 = Cons(Make_Symbol("define"), begin_subexpr1)
	begin_subexpr2 := SCMT_Nil
	begin_subexpr2 = Cons(Make_SCMT(3), begin_subexpr2)
	begin_subexpr2 = Cons(Make_Symbol("x"), begin_subexpr2)
	begin_subexpr2 = Cons(Make_Symbol("*"), begin_subexpr2)
	begin_expr := SCMT_Nil
	begin_expr = Cons(begin_subexpr2, begin_expr)
	begin_expr = Cons(begin_subexpr1, begin_expr)
	begin_expr = Cons(Make_Symbol("begin"), begin_expr)
	begin_result := begin_expr.Eval(env)
	if reflect.TypeOf(begin_result) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error(reflect.TypeOf(begin_result))
	}
	if begin_result.String() != "30" {
		t.Error(begin_result)
	}

	// Another complex one. We need to make this:
	// (let ((a 11)
	//       (b 12))
	//   (* a b))
	// Note: This should clobber the top-level environment's symbol "A"
	let_a_expr := SCMT_Nil
	let_a_expr = Cons(Make_SCMT(11), let_a_expr)
	let_a_expr = Cons(Make_Symbol("a"), let_a_expr)
	let_b_expr := SCMT_Nil
	let_b_expr = Cons(Make_SCMT(12), let_b_expr)
	let_b_expr = Cons(Make_Symbol("b"), let_b_expr)
	let_variables_expr := SCMT_Nil
	let_variables_expr = Cons(let_b_expr, let_variables_expr)
	let_variables_expr = Cons(let_a_expr, let_variables_expr)
	let_body_expr := SCMT_Nil
	let_body_expr = Cons(Make_Symbol("b"), let_body_expr)
	let_body_expr = Cons(Make_Symbol("a"), let_body_expr)
	let_body_expr = Cons(Make_Symbol("*"), let_body_expr)
	let_expr := SCMT_Nil
	let_expr = Cons(let_body_expr, let_expr)
	let_expr = Cons(let_variables_expr, let_expr)
	let_expr = Cons(Make_Symbol("let"), let_expr)
	let_result := let_expr.Eval(env)
	if reflect.TypeOf(let_result) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error(reflect.TypeOf(let_result))
	}
	if let_result.String() != "132" {
		t.Error(let_result)
	}

	// lambda
}

func Test_Proc(t *testing.T) {
	env := EnvSimple()
	
	// We're testing this:
	// ((lambda (n) (* n n)) 123) => 15129
	
	// args = (n)
	args := SCMT_Nil
	args = Cons(Make_Symbol("n"), args)

	// body = ((* i j))
	// Note the nesting, because begin is implied in the body
	body_inner := SCMT_Nil
	body_inner = Cons(Make_Symbol("n"), body_inner)
	body_inner = Cons(Make_Symbol("n"), body_inner)
	body_inner = Cons(Make_Symbol("*"), body_inner)
	body := SCMT_Nil
	body = Cons(body_inner, body)
	
	proc := Make_Proc(args, body, env)
	
	expr := SCMT_Nil
	expr = Cons(Make_SCMT(123), expr)
	expr = Cons(proc, expr)

	result := expr.Eval(env)
	if reflect.TypeOf(result) != reflect.TypeOf(&SCMT_Integer{}) {
		t.Error(reflect.TypeOf(result))
	}
	if result.String() != "15129" {
		t.Error(result)
	}
}
