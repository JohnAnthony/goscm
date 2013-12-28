package main

// TODO:
// Fix whitespace being leftover after a list expression
// Comments
// quote '
// quasiquote `
// unquote ,
// scm_float
// Escaped characters
// Scm Functions:
// eq? equal? = 
// atom? list?
// zero? false? true? not
// or and
// apply
// begin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type ScmType int

const (
	scm_string ScmType = iota
	scm_int
	scm_symbol
	scm_pair
	scm_func
	scm_gofunc
	scm_emptylist
	scm_bool
)

type ScmPair struct {
	car *Cell
	cdr *Cell
}

type Cell struct {
	stype ScmType
	value interface{}
}

func SCMBool(b bool) *Cell {
	return &Cell {
		stype: scm_bool,
		value: &b,
	}
}

func NewEnvironment() *Cell {
	var env *Cell = EmptyList()
	env = AddRawGoFunc(env, "quote", scm_quote)
	env = AddRawGoFunc(env, "define", scm_define)
	env = AddRawGoFunc(env, "lambda", scm_lambda)
	env = AddRawGoFunc(env, "cons", scm_cons)
	env = AddRawGoFunc(env, "car", scm_car)
	env = AddRawGoFunc(env, "cdr", scm_cdr)
	env = AddRawGoFunc(env, "display", scm_display)
	env = AddRawGoFunc(env, "if", nil)
	env = AddRawGoFunc(env, "+", scm_add)
	env = AddRawGoFunc(env, "-", scm_subtract)
	env = AddRawGoFunc(env, "*", scm_multiplication)
	env = AddRawGoFunc(env, "/", scm_division)
	return env
}

// recursively duplicate scheme objects
func duplicate(cell *Cell) *Cell {
	if cell == nil {
		return nil
	}
	
	switch cell.stype {
	case scm_string:
		fallthrough
	case scm_emptylist:
		fallthrough
	case scm_int:
		fallthrough
	case scm_symbol:
		fallthrough
	case scm_bool:
		fallthrough
	case scm_gofunc:
		val := cell.value
		return &Cell {
			stype: cell.stype,
			value: val,
		}
	case scm_func:
		fallthrough
	case scm_pair:
		return &Cell {
			stype: cell.stype,
			value: &ScmPair { car: duplicate(car(cell)), cdr: duplicate(cdr(cell)) },
		}
	}

	// If we get here it's an error
	return nil
}

func eval(env *Cell, expr *Cell) *Cell {
	switch expr.stype {
	case scm_int:
		fallthrough
	case scm_string:
		fallthrough
	case scm_gofunc:
		return expr
	case scm_symbol:
		return symbolLookup(env, expr.value.(string))
	case scm_func:
		return expr
	case scm_emptylist:
		return expr
	case scm_bool:
		return expr
	case scm_pair:
		// otherwise, we're using golang functions
		funcsym := car(expr).value.(string)
		tail := duplicate(cdr(expr))
		f := symbolLookup(env, funcsym)

		// Symbol not found
		if f == nil {
			return nil
		}

		// If special case
		if funcsym == "if" {
			pred := eval(env, car(tail)).value.(*bool)
			fst := car(cdr(tail))
			snd := car(cdr(cdr(tail)))

			if *pred == true {
				return eval(env, fst)
			} else {
				return eval(env, snd)
			}
		}

		// We don't eval if quoting
		if funcsym != "quote" && funcsym != "lambda" {
			var e *Cell

			// Special case for define - we don't eval the first symbol
			if funcsym == "define" {
				e = cdr(tail)
			} else {
				e = tail
			}

			// Not a special case - eval everything and build a new "tail"
			for ; e.stype != scm_emptylist; e = cdr(e) {
				e.value.(*ScmPair).car = eval(env, car(e))
			}
		}

		// This is everything we need for scheme functions
		if f.stype == scm_func {
			// Zip our symbols and values into a new environment
			subenv := env

			for symb, val := car(f), tail
			symb.stype != scm_emptylist && val.stype != scm_emptylist
			symb, val = cdr(symb), cdr(val) {
				pair := cons(car(symb), car(val))
				subenv = cons(pair, subenv)
			}

			var ret *Cell
			for sub := cdr(f); sub.stype != scm_emptylist; sub = cdr(sub) {
				ret = eval(subenv, car(sub))
			}

			return ret
		}

		// If we reach this point we should be handling scm_gofunc types

		ret := f.value.(func(*Cell) *Cell)(tail)

		// Handle define
		if funcsym == "define" {
			// Shuffle our new define into the place of the old env head
			tmpenv := *env
			*env = *cons(ret, &tmpenv)

			return nil
		}

		// Function is not define, return normally
		return ret
	}

	// Getting here is an error
	return nil
}

func AddRawGoFunc(env *Cell, symb string, f func(*Cell) *Cell) *Cell {
	newcar := &Cell{
		stype: scm_symbol,
		value: symb,
	}
	newcdr := &Cell{
		stype: scm_gofunc,
		value: f,
	}
	newcell := cons(newcar, newcdr)
	return cons(newcell, env)
}

func symbolLookup(env *Cell, symb string) *Cell {
	if env.stype == scm_emptylist {
		return nil
	}

	for e := env; e.stype != scm_emptylist; e = e.value.(*ScmPair).cdr {
		cell := e.value.(*ScmPair)
		symbpair := cell.car.value.(*ScmPair)
		label := symbpair.car.value.(string)
		f := symbpair.cdr
		if label == symb {
			return f
		}
	}

	// Couldn't find symbol
	return nil
}

func EmptyList() *Cell {
	return &Cell {
		stype: scm_emptylist,
	}
}

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	}
	return false
}

func chompspace(in *bufio.Reader) {
	var c byte
	for {
		c, _ = in.ReadByte()
		if !isSpace(c) {
			break
		}
	}
	in.UnreadByte()
}

func getexpr(in *bufio.Reader) *Cell {
	chompspace(in)
	symbol := make([]byte, 0)
	c, err := in.ReadByte()
	if err != nil {
		in.UnreadByte()
		return nil
	}

	// After chomping a '\n' means we're sticking on EOF
	if c == '\n' {
		return nil
	}

	// Scm Pair
	if c == '(' {
		var head *Cell = nil
		nexp := getexpr(in)
		if nexp.stype == scm_emptylist {
			return EmptyList()
		}
		
		head = &Cell{
			stype: scm_pair,
			value: &ScmPair{car: nexp},
		}
		tip := head

		for {
			nexp = getexpr(in)
			if nexp.stype == scm_emptylist {
				tip.value.(*ScmPair).cdr = EmptyList()
				break
			}
			tip.value.(*ScmPair).cdr = &Cell{
				stype: scm_pair,
				value: &ScmPair{car: nexp},
			}
			tip = cdr(tip)
		}
		return head
	}

	// Scm Pair Close
	if c == ')' {
		return EmptyList()
	}

	// Handle strings
	if c == '"' {
		for c, _ = in.ReadByte(); c != '"'; c, _ = in.ReadByte() {
			symbol = append(symbol, c)
		}
		return &Cell{
			stype: scm_string,
			value: string(symbol),
		}
	}

	for !isSpace(c) && c != ')' {
		symbol = append(symbol, c)
		c, _ = in.ReadByte()
	}

	if c == ')' {
		in.UnreadByte()
	}

	//Booleans
	if string(symbol) == "#t" {
		return SCMBool(true)
	}
	if string(symbol) == "#f" {
		return SCMBool(false)
	}

	// Integers
	if symbol[0] >= '0' && symbol[0] <= '9' {
		n, err := strconv.Atoi(string(symbol))
		if err != nil {
			fmt.Println(err)
		}

		return &Cell{
			stype: scm_int,
			value: &n,
		}
	}

	// Symbol
	return &Cell{
		stype: scm_symbol,
		value: string(symbol),
	}
}

func display(expr *Cell) {
	if expr == nil {
		fmt.Printf("nil")
		return
	}

	switch expr.stype {
	case scm_int:
		fmt.Printf("%d", *expr.value.(*int))
	case scm_string:
		fmt.Printf("\"%s\"", expr.value.(string))
	case scm_symbol:
		fmt.Printf("#<symbol %s>", expr.value.(string))
	case scm_emptylist:
		fmt.Printf("'()")
	case scm_bool:
		tf := *expr.value.(*bool)
		if tf == true {
			fmt.Printf("#t")
		} else {
			fmt.Printf("#f")
		}
	case scm_pair:
		// This is an ACTUAL pair
		if cdr(expr).stype != scm_emptylist && cdr(expr).stype != scm_pair {
			fmt.Printf("(")
			display(car(expr))
			fmt.Printf(" . ")
			display(cdr(expr))
			fmt.Printf(")")
			return
		}

		// This is a list
		fmt.Printf("(")
		for e := expr; e.stype != scm_emptylist; e = cdr(e) {
			display(car(e))
			if cdr(e).stype != scm_emptylist {
				fmt.Printf(" ")
			}
		}
		fmt.Printf(")")
	case scm_func:
		fmt.Printf("#<func>")
	case scm_gofunc:
		fmt.Printf("#<gofunc>")
	default:
		fmt.Printf("#<error>")
	}
}

// Go versions of scm functions

func cons(a *Cell, b *Cell) *Cell {
	return &Cell{
		stype: scm_pair,
		value: &ScmPair{car: a, cdr: b},
	}
}

func car(lst *Cell) *Cell {
	return lst.value.(*ScmPair).car
}

func cdr(lst *Cell) *Cell {
	return lst.value.(*ScmPair).cdr
}

// Scm functions

func scm_quote(tail *Cell) *Cell {
	return car(tail)
}

func scm_define(tail *Cell) *Cell {
	return &Cell{
		stype: scm_pair,
		value: &ScmPair{car: car(tail), cdr: car(cdr(tail))},
	}
}

func scm_lambda(tail *Cell) *Cell {
	return &Cell{
		stype: scm_func,
		value: &ScmPair{car: car(tail), cdr: cdr(tail)},
	}
}

func scm_cons(tail *Cell) *Cell {
	a := car(tail)
	b := car(cdr(tail))
	return cons(a, b)
}

func scm_car(tail *Cell) *Cell {
	return car(car(tail))
}

func scm_cdr(tail *Cell) *Cell {
	return cdr(car(tail))
}

func scm_display(tail *Cell) *Cell {
	display(car(tail))
	fmt.Println("")
	return nil
}

func scm_add(tail *Cell) *Cell {
	ret := 0

	for e := tail; e.stype != scm_emptylist; e = cdr(e) {
		ret += *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_subtract(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e.stype != scm_emptylist; e = cdr(e) {
		ret -= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_multiplication(tail *Cell) *Cell {
	ret := 1

	for e := tail; e.stype != scm_emptylist; e = cdr(e) {
		ret *= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_division(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e.stype != scm_emptylist; e = cdr(e) {
		ret /= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := NewEnvironment()
	
	for {
		// Read
		expr := getexpr(reader)
		if expr == nil {
			break
		}
		// Eval
		ret := eval(env, expr)
		// Print
		fmt.Printf("|> ")
		display(ret)
		fmt.Println("")
	}
}
