package simple

import (
	"testing"
	"reflect"
	"github.com/JohnAnthony/goscm"
)

func Test_EnvSimple(t *testing.T) {
	env := EnvSimple()

	// Test define
	goscm.Make_List(
		goscm.Make_Symbol("define"),
		goscm.Make_Symbol("a"),
		goscm.Make_SCMT(1234),
	).Eval(env)
	a_result := goscm.Make_Symbol("a").Eval(env)
	if reflect.TypeOf(a_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(a_result))
	}
	if a_result.String() != "1234" {
		t.Error(a_result)
	}
	
	// Test +
	add_result := goscm.Make_List(
		goscm.Make_Symbol("+"),
		goscm.Make_SCMT(321),
		goscm.Make_Symbol("a"),
		goscm.Make_SCMT(222),
		goscm.Make_SCMT(1000),
	).Eval(env)
	if reflect.TypeOf(add_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(add_result))
	}
	if add_result.String() != "2777" {
		t.Error(add_result)
	}
	
	// Test -
	sub_result := goscm.Make_List(
		goscm.Make_Symbol("-"),
		goscm.Make_SCMT(100),
		goscm.Make_SCMT(9),
		goscm.Make_SCMT(90),
	).Eval(env)
	if reflect.TypeOf(sub_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(sub_result))
	}
	if sub_result.String() != "1" {
		t.Error(sub_result)
	}

	// Test *
	mult_result := goscm.Make_List(
		goscm.Make_Symbol("*"),
		goscm.Make_SCMT(1),
		goscm.Make_SCMT(2),
		goscm.Make_SCMT(3),
		goscm.Make_SCMT(4),
		goscm.Make_SCMT(5),
	).Eval(env)
	if reflect.TypeOf(mult_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(mult_result))
	}
	if mult_result.String() != "120" {
		t.Error(mult_result)
	}
	
	// Test /
	div_result := goscm.Make_List(
		goscm.Make_Symbol("/"),
		goscm.Make_SCMT(66),
		goscm.Make_SCMT(3),
		goscm.Make_SCMT(2),
	).Eval(env)
	if reflect.TypeOf(div_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(div_result))
	}
	if div_result.String() != "11" {
		t.Error(div_result)
	}
	
	// Test cons
	// TODO
	
	// Test car
	car_result := goscm.Make_List(
		goscm.Make_Symbol("car"),
		goscm.Make_List(
			goscm.Make_Symbol("cons"),
			goscm.Make_SCMT("FOO"),
			goscm.Make_SCMT("BAR"),
		),
	).Eval(env)
	if reflect.TypeOf(car_result) != reflect.TypeOf(&goscm.SCMT_String{}) {
		t.Error(reflect.TypeOf(car_result))
	}
	if car_result.String() != "\"FOO\"" {
		t.Error(car_result)
	}

	// Building this list: (cdr (cons "FOO" "BAR"))
	cdr_expr := goscm.SCMT_Nil
	cdr_expr = goscm.Cons(goscm.Make_SCMT("BAR"), cdr_expr)
	cdr_expr = goscm.Cons(goscm.Make_SCMT("FOO"), cdr_expr)
	cdr_expr = goscm.Cons(goscm.Make_Symbol("cons"), cdr_expr)
	cdr_expr = goscm.Cons(cdr_expr, goscm.SCMT_Nil)
	cdr_expr = goscm.Cons(goscm.Make_Symbol("cdr"), cdr_expr)
	cdr_result := cdr_expr.Eval(env)
	if reflect.TypeOf(cdr_result) != reflect.TypeOf(&goscm.SCMT_String{}) {
		t.Error(reflect.TypeOf(cdr_result))
	}
	if cdr_result.String() != "\"BAR\"" {
		t.Error(cdr_result)
	}
	
	cons_expr := goscm.SCMT_Nil
	cons_expr = goscm.Cons(goscm.Make_SCMT(5), cons_expr)
	cons_expr = goscm.Cons(goscm.Make_SCMT(2), cons_expr)
	cons_expr = goscm.Cons(goscm.Make_Symbol("cons"), cons_expr)
	cons_result := cons_expr.Eval(env)
	if reflect.TypeOf(cons_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(cons_result))
	}
	if cons_result.String() != "(2 . 5)" {
		t.Error(cons_result)
	}
	
	// This is a little complex - we're actually making a list whose first
	// element is the symbol "quote", whose second element is a complicated list
	// and which is then terminated by goscm.SCMT_Nil.
	quote_expr := goscm.SCMT_Nil
	quote_expr = goscm.Cons(goscm.Make_SCMT(9), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(8), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(7), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_Symbol("honk"), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(6), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(5), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(4), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(3), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_Symbol("squeak"), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(2), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_SCMT(1), quote_expr)
	quote_expr = goscm.Cons(goscm.Make_Symbol("roar"), quote_expr)
	quote_expr = goscm.Cons(quote_expr, goscm.SCMT_Nil)
	quote_expr = goscm.Cons(goscm.Make_Symbol("quote"), quote_expr)
	quote_result := quote_expr.Eval(env)
	if reflect.TypeOf(quote_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
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
	begin_subexpr1 := goscm.SCMT_Nil
	begin_subexpr1 = goscm.Cons(goscm.Make_SCMT(10), begin_subexpr1)
	begin_subexpr1 = goscm.Cons(goscm.Make_Symbol("x"), begin_subexpr1)
	begin_subexpr1 = goscm.Cons(goscm.Make_Symbol("define"), begin_subexpr1)
	begin_subexpr2 := goscm.SCMT_Nil
	begin_subexpr2 = goscm.Cons(goscm.Make_SCMT(3), begin_subexpr2)
	begin_subexpr2 = goscm.Cons(goscm.Make_Symbol("x"), begin_subexpr2)
	begin_subexpr2 = goscm.Cons(goscm.Make_Symbol("*"), begin_subexpr2)
	begin_expr := goscm.SCMT_Nil
	begin_expr = goscm.Cons(begin_subexpr2, begin_expr)
	begin_expr = goscm.Cons(begin_subexpr1, begin_expr)
	begin_expr = goscm.Cons(goscm.Make_Symbol("begin"), begin_expr)
	begin_result := begin_expr.Eval(env)
	if reflect.TypeOf(begin_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
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
	let_a_expr := goscm.SCMT_Nil
	let_a_expr = goscm.Cons(goscm.Make_SCMT(11), let_a_expr)
	let_a_expr = goscm.Cons(goscm.Make_Symbol("a"), let_a_expr)
	let_b_expr := goscm.SCMT_Nil
	let_b_expr = goscm.Cons(goscm.Make_SCMT(12), let_b_expr)
	let_b_expr = goscm.Cons(goscm.Make_Symbol("b"), let_b_expr)
	let_variables_expr := goscm.SCMT_Nil
	let_variables_expr = goscm.Cons(let_b_expr, let_variables_expr)
	let_variables_expr = goscm.Cons(let_a_expr, let_variables_expr)
	let_body_expr := goscm.SCMT_Nil
	let_body_expr = goscm.Cons(goscm.Make_Symbol("b"), let_body_expr)
	let_body_expr = goscm.Cons(goscm.Make_Symbol("a"), let_body_expr)
	let_body_expr = goscm.Cons(goscm.Make_Symbol("*"), let_body_expr)
	let_expr := goscm.SCMT_Nil
	let_expr = goscm.Cons(let_body_expr, let_expr)
	let_expr = goscm.Cons(let_variables_expr, let_expr)
	let_expr = goscm.Cons(goscm.Make_Symbol("let"), let_expr)
	let_result := let_expr.Eval(env)
	if reflect.TypeOf(let_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(let_result))
	}
	if let_result.String() != "132" {
		t.Error(let_result)
	}

	// Test with this expression:
	// ((lambda (x) (+ x 2)) 58835) => 58837
	lambda_args := goscm.Cons(goscm.Make_Symbol("x"), goscm.SCMT_Nil)
	lambda_body := goscm.SCMT_Nil
	lambda_body = goscm.Cons(goscm.Make_SCMT(2), lambda_body)
	lambda_body = goscm.Cons(goscm.Make_Symbol("x"), lambda_body)
	lambda_body = goscm.Cons(goscm.Make_Symbol("+"), lambda_body)
	lambda_expr := goscm.SCMT_Nil
	lambda_expr = goscm.Cons(lambda_body, lambda_expr)
	lambda_expr = goscm.Cons(lambda_args, lambda_expr)
	lambda_expr = goscm.Cons(goscm.Make_Symbol("lambda"), lambda_expr)
	full_expr := goscm.SCMT_Nil
	full_expr = goscm.Cons(goscm.Make_SCMT(58835), full_expr)
	full_expr = goscm.Cons(lambda_expr, full_expr)
	lambda_result := full_expr.Eval(env)
	if reflect.TypeOf(lambda_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(lambda_result))
	}
	if lambda_result.String() != "58837" {
		t.Error(lambda_result)
	}

	// Test with these expressions:
	// (define square (lambda (x) (* x x)))
	// (map square (quote (2 3 4 5 6))) => (4 9 16 25 36)
	squ_lambda_body := goscm.SCMT_Nil
	squ_lambda_body = goscm.Cons(goscm.Make_Symbol("x"), squ_lambda_body)
	squ_lambda_body = goscm.Cons(goscm.Make_Symbol("x"), squ_lambda_body)
	squ_lambda_body = goscm.Cons(goscm.Make_Symbol("*"), squ_lambda_body)
	squ_lambda := goscm.SCMT_Nil
	squ_lambda = goscm.Cons(squ_lambda_body, squ_lambda)
	squ_lambda = goscm.Cons(goscm.Cons(goscm.Make_Symbol("x"), goscm.SCMT_Nil), squ_lambda)
	squ_lambda = goscm.Cons(goscm.Make_Symbol("lambda"), squ_lambda)
	squ_define := goscm.SCMT_Nil
	squ_define = goscm.Cons(squ_lambda, squ_define)
	squ_define = goscm.Cons(goscm.Make_Symbol("square"), squ_define)
	squ_define = goscm.Cons(goscm.Make_Symbol("define"), squ_define)
	squ_define.Eval(env)

	map_expr := goscm.SCMT_Nil
	map_expr = goscm.Cons(goscm.Make_SCMT(6), map_expr)
	map_expr = goscm.Cons(goscm.Make_SCMT(5), map_expr)
	map_expr = goscm.Cons(goscm.Make_SCMT(4), map_expr)
	map_expr = goscm.Cons(goscm.Make_SCMT(3), map_expr)
	map_expr = goscm.Cons(goscm.Make_SCMT(2), map_expr)
	map_expr = goscm.Cons(map_expr, goscm.SCMT_Nil)
	map_expr = goscm.Cons(goscm.Make_Symbol("quote"), map_expr)
	map_expr = goscm.Cons(map_expr, goscm.SCMT_Nil)
	map_expr = goscm.Cons(goscm.Make_Symbol("square"), map_expr)
	map_expr = goscm.Cons(goscm.Make_Symbol("map"), map_expr)
	map_result := map_expr.Eval(env)
	if reflect.TypeOf(map_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(map_result))
	}
	if map_result.String() != "(4 9 16 25 36)" {
		t.Error(map_result)
	}
	
	// Test apply with this expression
	// (apply + (quote (1 2 3 4 5 6))) => 21
	apply_expr := goscm.SCMT_Nil
	apply_expr = goscm.Cons(goscm.Make_SCMT(6), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_SCMT(5), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_SCMT(4), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_SCMT(3), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_SCMT(2), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_SCMT(1), apply_expr)
	apply_expr = goscm.Cons(apply_expr, goscm.SCMT_Nil)
	apply_expr = goscm.Cons(goscm.Make_Symbol("quote"), apply_expr)
	apply_expr = goscm.Cons(apply_expr, goscm.SCMT_Nil)
	apply_expr = goscm.Cons(goscm.Make_Symbol("+"), apply_expr)
	apply_expr = goscm.Cons(goscm.Make_Symbol("apply"), apply_expr)
	apply_result := apply_expr.Eval(env)
	if reflect.TypeOf(apply_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(apply_result))
	}
	if apply_result.String() != "21" {
		t.Error(apply_result)
	}
}
