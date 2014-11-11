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
	a_result, err := goscm.Make_Symbol("a").Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(a_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(a_result))
	}
	if a_result.String() != "1234" {
		t.Error(a_result)
	}
	
	// Test +
	add_result, err := goscm.Make_List(
		goscm.Make_Symbol("+"),
		goscm.Make_SCMT(321),
		goscm.Make_Symbol("a"),
		goscm.Make_SCMT(222),
		goscm.Make_SCMT(1000),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(add_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(add_result))
	}
	if add_result.String() != "2777" {
		t.Error(add_result)
	}
	
	// Test -
	sub_result, err := goscm.Make_List(
		goscm.Make_Symbol("-"),
		goscm.Make_SCMT(100),
		goscm.Make_SCMT(9),
		goscm.Make_SCMT(90),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(sub_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(sub_result))
	}
	if sub_result.String() != "1" {
		t.Error(sub_result)
	}

	// Test *
	mult_result, err := goscm.Make_List(
		goscm.Make_Symbol("*"),
		goscm.Make_SCMT(1),
		goscm.Make_SCMT(2),
		goscm.Make_SCMT(3),
		goscm.Make_SCMT(4),
		goscm.Make_SCMT(5),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(mult_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(mult_result))
	}
	if mult_result.String() != "120" {
		t.Error(mult_result)
	}
	
	// Test /
	div_result, err := goscm.Make_List(
		goscm.Make_Symbol("/"),
		goscm.Make_SCMT(66),
		goscm.Make_SCMT(3),
		goscm.Make_SCMT(2),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(div_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(div_result))
	}
	if div_result.String() != "11" {
		t.Error(div_result)
	}
	
	// Test cons
	cons_result, err := goscm.Make_List(
		goscm.Make_Symbol("cons"),
		goscm.Make_SCMT(2),
		goscm.Make_SCMT(5),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(cons_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(cons_result))
	}
	if cons_result.String() != "(2 . 5)" {
		t.Error(cons_result)
	}
	
	// Test car
	car_result, err := goscm.Make_List(
		goscm.Make_Symbol("car"),
		goscm.Make_List(
			goscm.Make_Symbol("cons"),
			goscm.Make_SCMT(333),
			goscm.Make_SCMT(444),
		),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(car_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(car_result))
	}
	if car_result.String() != "333" {
		t.Error(car_result)
	}

	// Test cdr
	cdr_result, err := goscm.Make_List(
		goscm.Make_Symbol("cdr"),
		goscm.Make_List(
			goscm.Make_Symbol("cons"),
			goscm.Make_SCMT(555),
			goscm.Make_SCMT(666),
		),
	).Eval(env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(cdr_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(cdr_result))
	}
	if cdr_result.String() != "666" {
		t.Error(cdr_result)
	}
	
	// Test quote
	quote_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}
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
	begin_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}
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
	let_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(let_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(let_result))
	}
	if let_result.String() != "132" {
		t.Error(let_result)
	}

	// Test lambda
	// ((lambda (x) (+ x 2)) 58835) => 58837
	lambda_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}
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

	map_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(map_result) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(map_result))
	}
	if map_result.String() != "(4 9 16 25 36)" {
		t.Error(map_result)
	}
	
	// Test apply
	// (apply + (quote (1 2 3 4 5 6))) => 21
	apply_result, err := goscm.Make_List(
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
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(apply_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(apply_result))
	}
	if apply_result.String() != "21" {
		t.Error(apply_result)
	}
	
	// Test set!
	// (define x 10)
	// (set! x 33)
	// x => 33
	_, err = goscm.EvalStr("(define macamaga 10)", Read, env)
	if err != nil {
		t.Error(err)
	}
	_, err = goscm.EvalStr("(set! macamaga 33)", Read, env)
	if err != nil {
		t.Error(err)
	}
	set_result, err := goscm.EvalStr("macamaga", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(set_result) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(set_result))
	}
	if set_result.String() != "33" {
		t.Error(set_result)
	}
	
	// Test the alternative function definition syntax
	// (define (square-plus-cube n m)
    //   (+ (* n n)
	//      (* m m m)))
	// (square-plus-cube 11 17) => 5034
	_ , err = goscm.EvalStr(`(define (square-plus-cube n m)
                                (+ (* n n)
                                   (* m m m)))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	altdefine, err := goscm.EvalStr("(square-plus-cube 11 17)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(altdefine) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(altdefine))
	}
	if altdefine.String() != "5034" {
		t.Error(altdefine)
	}
	
	// Test the numerical equality operator for success
	// (= 182736 182736) => #t
	numeqt, err := goscm.EvalStr("(= 182736 182736 182736)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(numeqt) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(numeqt))
	}
	if numeqt.String() != "#t" {
		t.Error(numeqt)
	}
	
	// Test the numerical equality operator for success
	// (= 11 12) => #f
	numeqf, err := goscm.EvalStr("(= 11 12)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(numeqf) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(numeqf))
	}
	if numeqf.String() != "#f" {
		t.Error(numeqf)
	}
	
	// Test if special form for truth side
	// (if (= 10 10)
    //   (+ 2 3)
	//   (* 20 30)) => 5
	ift, err := goscm.EvalStr("(if (= 10 10) (+ 2 3) (* 20 30))", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ift) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(ift))
	}
	if ift.String() != "5" {
		t.Error(ift)
	}

	// Test if special form for false side
	// (if (= 10 11)
    //   (+ 2 3)
	//   (* 20 30)) => 600
	iff, err := goscm.EvalStr("(if (= 10 11) (+ 2 3) (* 20 30))", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(iff) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(iff))
	}
	if iff.String() != "600" {
		t.Error(iff)
	}
	
	// Check fib works just because we can
	// (define (fib n)
	//   (if (= n 2)
	//       1
	//       (if (= n 1)
	//           1
	//           (+ (fib (- n 1)) (fib (- n 2))))))
	// (fib 10) => 
/*	_, err = goscm.EvalStr(`(define (fib n)
	   (if (= n 2)
	       1
	       (if (= n 1)
	           1
	           (+ (fib (- n 1)) (fib (- n 2))))))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	fib, err := goscm.EvalStr("(fib 10)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(fib) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(fib))
	}
	if fib.String() != "55" {
		t.Error(fib)
	} */
	
	// Check reading with a dot in the notation
	dotcdr, err := goscm.EvalStr("(cdr '(123 . 456))", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(dotcdr) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(dotcdr))
	}
	if dotcdr.String() != "456" {
		t.Error(dotcdr)
	}
	
	// Check quote syntax
	// (map square '(2 3 4 5 6)) => (4 9 16 25 36
	quotesyn, err := goscm.EvalStr("(map square '(2 3 4 5 6))", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(quotesyn) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(quotesyn))
	}
	if quotesyn.String() != "(4 9 16 25 36)" {
		t.Error(quotesyn)
	}
	
	// Check that dotted function arguments work
	// (define (pairup a . rest)
	//   (map (lambda (b) (cons a b)) rest))
	// (pairup 3 4 5 6) => ((3 . 4) (3 . 5) (3 . 6))
	_, err = goscm.EvalStr(`(define (pairup a . rest)
                              (map (lambda (b) (cons a b)) rest))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	dotlam, err := goscm.EvalStr("(pairup 3 4 5 6)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(dotlam) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(dotlam))
	}
	if dotlam.String() != "((3 . 4) (3 . 5) (3 . 6))" {
		t.Error(dotlam)
	}

	// The same as above but with different syntax
	// (define pairup2
	//   (lambda (a . rest)
	//     (map (lambda (b) (cons a b)) rest)))
	// (pairup2 3 4 5 6) => ((3 . 4) (3 . 5) (3 . 6))
	_, err = goscm.EvalStr(`(define pairup2
                              (lambda (a . rest)
                                (map (lambda (b) (cons a b)) rest)))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	dotlam2, err := goscm.EvalStr("(pairup2 3 4 5 6)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(dotlam2) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(dotlam2))
	}
	if dotlam2.String() != "((3 . 4) (3 . 5) (3 . 6))" {
		t.Error(dotlam2)
	}

	// 'not' should return #t when passed #f
	nota, err := goscm.EvalStr("(not #f)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(nota) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(nota))
	}
	if nota.String() != "#t" {
		t.Error(nota)
	}

	// 'not' should return #f when passed an integer or anything but #t
	notb, err := goscm.EvalStr("(map not '(2364 #f #t 4))", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(notb) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(notb))
	}
	if notb.String() != "(#f #t #f #f)" {
		t.Error(notb)
	}

	// Test for < of a long list in which it is true
	ltt, err := goscm.EvalStr(`(and (< 2 100 101 200 203 1000)
                                    (< 3 5)
                                    (< 8 9 11 12))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ltt) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(ltt))
	}
	if ltt.String() != "#t" {
		t.Error(ltt)
	}

	// Test for < of a long list in which it is false
	ltf, err := goscm.EvalStr(`(or (< 2 100 100 1 203 1000)
                                   (< 4 3)
                                   (< 99 100 99 101))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ltf) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(ltf))
	}
	if ltf.String() != "#f" {
		t.Error(ltf)
	}

	// Test for > of a long list in which it is true
	gtt, err := goscm.EvalStr(`(and (> 1000 202 201 200 50 5 1)
                                    (> 9 8 7 6 5 4 3 2 1)
                                    (> 3 2))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(gtt) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(gtt))
	}
	if gtt.String() != "#t" {
		t.Error(gtt)
	}

	// Test for > of a long list in which it is false
	gtf, err := goscm.EvalStr(`(or (> 10000 5000 5000 20 30 0 44)
                                   (> 10 11)
                                   (> 100 99 100 99))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(gtf) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(gtf))
	}
	if gtf.String() != "#f" {
		t.Error(gtf)
	}

	// Test for 'and' when it should succeed
	andt, err := goscm.EvalStr("(and #t #t #t #t #t #t #t #t #t #t)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(andt) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(andt))
	}
	if andt.String() != "#t" {
		t.Error(andt)
	}

	// Test for 'and' when it should fail
	// TODO: Ensure that eval only happens if previous cases are true
	andf, err := goscm.EvalStr("(and #t #t #t #t #t #t #t #f #t #t)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(andf) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(andf))
	}
	if andf.String() != "#f" {
		t.Error(andf)
	}

	// Test for 'or' when it should succeed
	// TODO test that eval is only happening when it should
	ort, err := goscm.EvalStr("(or #f #f #f #f #f #f #t #f #f)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(ort) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(ort))
	}
	if ort.String() != "#t" {
		t.Error(ort)
	}

	// Test for 'or' when it should fail
	// TODO test that eval is only happening when it should
	orf, err := goscm.EvalStr("(or #f #f #f #f #f #f #f #f #f)", Read, env)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(orf) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(orf))
	}
	if orf.String() != "#f" {
		t.Error(orf)
	}

	// Test for '>=' when it should pass
	get, err := goscm.EvalStr(`(and (>= 10 10 10 10 9 9 9 9 8 8 7 7)
                                    (>= 11 11)
                                    (>= 9 8))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if get.String() != "#t" {
		t.Error(get)
	}
	
	// Test for '>=' when it should fail
	gef, err := goscm.EvalStr(`(or (>= 10 10 10 10 9 9 9 9 10 8 7)
                                   (>= 11 12))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if gef.String() != "#f" {
		t.Error(gef)
	}
		
	// Test for '<=' when it should pass
	let, err := goscm.EvalStr(`(and (<= 10 10 11 11 11 2001)
                                    (<= 11 12))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if let.String() != "#t" {
		t.Error(let)
	}
	
	// Test for '<=' when it should fail
	lef, err := goscm.EvalStr(`(or (<= 10 10 10 10 9 10 10 11 14)
                                   (<= 11 10))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if lef.String() != "#f" {
		t.Error(lef)
	}
	
	// Test for 'reverse'
	rev, err := goscm.EvalStr(`(reverse '(1 2 3 (+ 88 22) 4 6))`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if rev.String() != "(6 4 (+ 88 22) 3 2 1)" {
		t.Error(rev)
	}
	
	// Test for 'list'
	list_func, err := goscm.EvalStr(`(list 2 3 4 'a 3 5 67)`, Read, env)
	if err != nil {
		t.Error(err)
	}
	if list_func.String() != "(2 3 4 A 3 5 67)" {
		t.Error(list_func)
	}
}

func Test_Read(t *testing.T) {
	symbol, serr := ReadStr("anaconda")
	if serr != nil {
		t.Error(serr)
	}
	if reflect.TypeOf(symbol) != reflect.TypeOf(&goscm.SCMT_Symbol{}) {
		t.Error(reflect.TypeOf(symbol))
	}
	if symbol.String() != "ANACONDA" {
		t.Error(symbol)
	}

	int, ierr := ReadStr("1337")
	if ierr != nil {
		t.Error(ierr)
	}
	if reflect.TypeOf(int) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
		t.Error(reflect.TypeOf(int))
	}
	if int.String() != "1337" {
		t.Error(int)
	}

	list, lerr := ReadStr("(+ 111 222 333)")
	if lerr != nil {
		t.Error(lerr)
	}
	if reflect.TypeOf(list) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		t.Error(reflect.TypeOf(list))
	}
	if list.String() != "(+ 111 222 333)" {
		t.Error(list)
	}
	
	btrue, err := ReadStr("#t")
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(btrue) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(btrue))
	}
	if btrue.String() != "#t" {
		t.Error(btrue)
	}

	bfalse, err := ReadStr("#f")
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(bfalse) != reflect.TypeOf(&goscm.SCMT_Bool{}) {
		t.Error(reflect.TypeOf(bfalse))
	}
	if bfalse.String() != "#f" {
		t.Error(bfalse)
	}
}
