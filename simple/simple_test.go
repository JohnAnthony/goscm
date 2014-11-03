package simple

import (
	"testing"
	"reflect"
	"github.com/JohnAnthony/goscm"
)

func Test_Env(t *testing.T) {
	env := Env()

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
	cons_result := goscm.Make_List(
		goscm.Make_Symbol("cons"),
		goscm.Make_SCMT(2),
		goscm.Make_SCMT(5),
	).Eval(env)
	if reflect.TypeOf(cons_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(cons_result))
	}
	if cons_result.String() != "(2 . 5)" {
		t.Error(cons_result)
	}
	
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

	// Test cdr
	cdr_result := goscm.Make_List(
		goscm.Make_Symbol("cdr"),
		goscm.Make_List(
			goscm.Make_Symbol("cons"),
			goscm.Make_SCMT("FOO"),
			goscm.Make_SCMT("BAR"),
		),
	).Eval(env)
	if reflect.TypeOf(cdr_result) != reflect.TypeOf(&goscm.SCMT_String{}) {
		t.Error(reflect.TypeOf(cdr_result))
	}
	if cdr_result.String() != "\"BAR\"" {
		t.Error(cdr_result)
	}
	
	// Test quote
	quote_result := goscm.Make_List(
		goscm.Make_Symbol("quote"),
		goscm.Make_List(
			goscm.Make_Symbol("roar"),
			goscm.Make_SCMT(1),
			goscm.Make_SCMT(2),
			goscm.Make_Symbol("squeak"),
			goscm.Make_SCMT(3),
			goscm.Make_SCMT(4),
			goscm.Make_SCMT(5),
			goscm.Make_SCMT(6),
			goscm.Make_Symbol("honk"),
			goscm.Make_SCMT(7),
			goscm.Make_SCMT(8),
			goscm.Make_SCMT(9),
		),
	).Eval(env)
	if reflect.TypeOf(quote_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(quote_result))
	}
	if quote_result.String() != "(ROAR 1 2 SQUEAK 3 4 5 6 HONK 7 8 9)" {
		t.Error(quote_result)
	}

	// Test begin
	// (begin
	//  (define x 10)
	//  (* x 3)) => 30
	begin_result := goscm.Make_List(
		goscm.Make_Symbol("begin"),
		goscm.Make_List(
			goscm.Make_Symbol("define"),
			goscm.Make_Symbol("x"),
			goscm.Make_SCMT(10),
		),
		goscm.Make_List(
			goscm.Make_Symbol("*"),
			goscm.Make_Symbol("x"),
			goscm.Make_SCMT(3),
		),
	).Eval(env)
	if reflect.TypeOf(begin_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(begin_result))
	}
	if begin_result.String() != "30" {
		t.Error(begin_result)
	}

	// Test let
	// (let ((a 11)
	//       (b 12))
	//   (* a b))
	// Note: This should clobber the top-level environment's symbol "A"
	let_result := goscm.Make_List(
		goscm.Make_Symbol("let"),
		goscm.Make_List(
			goscm.Make_List(
				goscm.Make_Symbol("a"),
				goscm.Make_SCMT(11),
			),
			goscm.Make_List(
				goscm.Make_Symbol("b"),
				goscm.Make_SCMT(12),
			),
		),
		goscm.Make_List(
			goscm.Make_Symbol("*"),
			goscm.Make_Symbol("a"),
			goscm.Make_Symbol("b"),
		),
	).Eval(env)
	if reflect.TypeOf(let_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(let_result))
	}
	if let_result.String() != "132" {
		t.Error(let_result)
	}

	// Test lambda
	// ((lambda (x) (+ x 2)) 58835) => 58837
	lambda_result := goscm.Make_List(
		goscm.Make_List(
			goscm.Make_Symbol("lambda"),
			goscm.Make_List(
				goscm.Make_Symbol("x"),
			),
			goscm.Make_List(
				goscm.Make_Symbol("+"),
				goscm.Make_Symbol("x"),
				goscm.Make_SCMT(2),
			),
		),
		goscm.Make_SCMT(58835),
	).Eval(env)
	if reflect.TypeOf(lambda_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(lambda_result))
	}
	if lambda_result.String() != "58837" {
		t.Error(lambda_result)
	}

	// Test map
	// (define square (lambda (x) (* x x)))
	// (map square (quote (2 3 4 5 6))) => (4 9 16 25 36)
	goscm.Make_List(
		goscm.Make_Symbol("define"),
		goscm.Make_Symbol("square"),
		goscm.Make_List(
			goscm.Make_Symbol("lambda"),
			goscm.Make_List(
				goscm.Make_Symbol("x"),
			),
			goscm.Make_List(
				goscm.Make_Symbol("*"),
				goscm.Make_Symbol("x"),
				goscm.Make_Symbol("x"),
			),
		),
	).Eval(env)

	map_result := goscm.Make_List(
		goscm.Make_Symbol("map"),
		goscm.Make_Symbol("square"),
		goscm.Make_List(
			goscm.Make_Symbol("quote"),
			goscm.Make_List(
				goscm.Make_SCMT(2),
				goscm.Make_SCMT(3),
				goscm.Make_SCMT(4),
				goscm.Make_SCMT(5),
				goscm.Make_SCMT(6),
			),
		),
	).Eval(env)

	if reflect.TypeOf(map_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(map_result))
	}
	if map_result.String() != "(4 9 16 25 36)" {
		t.Error(map_result)
	}
	
	// Test apply
	// (apply + (quote (1 2 3 4 5 6))) => 21
	apply_result := goscm.Make_List(
		goscm.Make_Symbol("apply"),
		goscm.Make_Symbol("+"),
		goscm.Make_List(
			goscm.Make_Symbol("quote"),
			goscm.Make_List(
				goscm.Make_SCMT(1),
				goscm.Make_SCMT(2),
				goscm.Make_SCMT(3),
				goscm.Make_SCMT(4),
				goscm.Make_SCMT(5),
				goscm.Make_SCMT(6),
			),
		),
	).Eval(env)
	if reflect.TypeOf(apply_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(apply_result))
	}
	if apply_result.String() != "21" {
		t.Error(apply_result)
	}
}

func Test_Read(t *testing.T) {
	symbol := Read("a")
	if reflect.TypeOf(symbol) != reflect.TypeOf(&goscm.SCMT_Symbol{}) {
		t.Error(reflect.TypeOf(symbol))
	}
	if symbol.String() != "A" {
		t.Error(symbol)
	}
}
