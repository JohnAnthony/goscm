package main

// TODO:
// define
// lambda
// Comments
// quote '
// quasiquote `
// unquote ,
// scm_bool (including #f and #t syntax reading)
// if
// or
// and
// scm_float

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
	scm_pairclose
	scm_gofunc
)

type ScmPair struct {
	car *Cell
	cdr *Cell
}

type Cell struct {
	stype ScmType
	value interface{}
}

func NewEnvironment() *Cell {
	// Add all of our special forms
	quotepair := cons(&Cell{stype: scm_symbol, value: "quote"}, &Cell{stype: scm_gofunc, value: scm_quote})
	newenv := cons(quotepair, nil)

	return newenv
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
	}

	// From here on is handling scm_pairs

	funcsym := car(expr).value.(string)
	tail := cdr(expr)
	f := symbolLookup(env, funcsym)

	if f == nil {
		return nil
	}

	// We don't eval if quoting!
	if funcsym != "quote" {
		// Eval all the tail args first
		for e := tail; e != nil; e = cdr(e) {
			e.value.(*ScmPair).car = eval(env, car(e))
		}
	}

	// Go source functions
	if f.stype == scm_gofunc {
		return f.value.(func(*Cell) *Cell)(tail)
	}

	// Error because we've got an unhandled type
	return nil
}

func envupdate(env *Cell, add *Cell) {
	shifted := *env
	*env = *cons(add, &shifted)
}

func AddRawGoFunc(env *Cell, symb string, f func(*Cell) *Cell) {
	newcar := &Cell{
		stype: scm_symbol,
		value: symb,
	}
	newcdr := &Cell{
		stype: scm_gofunc,
		value: f,
	}
	newcell := cons(newcar, newcdr)
	envupdate(env, newcell)
}

func symbolLookup(env *Cell, symb string) *Cell {
	if env == nil {
		return nil
	}

	for e := env; e != nil; e = e.value.(*ScmPair).cdr {
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

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	}
	return false
}

func getexpr(in *bufio.Reader) *Cell {
	var c byte
	symbol := make([]byte, 0)

	// Skip over whitespace
	for {
		c, _ = in.ReadByte()
		if !isSpace(c) {
			break
		}
	}

	// Scm Pair
	if c == '(' {
		var head *Cell = nil
		nexp := getexpr(in)
		if nexp == nil {
			return nil
		}

		head = &Cell{
			stype: scm_pair,
			value: &ScmPair{car: nexp},
		}
		tip := head
		for {
			nexp = getexpr(in)
			if nexp.stype == scm_pairclose {
				tip.value.(*ScmPair).cdr = nil
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
		return &Cell{
			stype: scm_pairclose,
		}
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
		fmt.Printf("%s", expr.value.(string))
	case scm_pair:
		// This is an ACTUAL pair
		if cdr(expr) != nil && cdr(expr).stype != scm_pair {
			fmt.Printf("(")
			display(car(expr))
			fmt.Printf(" . ")
			display(cdr(expr))
			fmt.Printf(")")
			return
		}
		
		// This is a list
		fmt.Printf("(")
		for e := expr; e != nil; e = cdr(e) {
			display(car(e))
			if cdr(e) != nil {
				fmt.Printf(" ")
			}
		}
		fmt.Printf(")")
	case scm_gofunc:
		fmt.Printf("<#gofunc>")
	default:
		fmt.Printf("<#error>")
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

func scm_add(tail *Cell) *Cell {
	ret := 0

	for e := tail; e != nil; e = cdr(e) {
		ret += *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_subtract(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e != nil; e = cdr(e) {
		ret -= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_multiplication(tail *Cell) *Cell {
	ret := 1

	for e := tail; e != nil; e = cdr(e) {
		ret *= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_division(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e != nil; e = cdr(e) {
		ret /= *car(e).value.(*int)
	}

	return &Cell{
		stype: scm_int,
		value: &ret,
	}
}

func scm_quote(tail *Cell) *Cell {
	return car(tail)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := NewEnvironment()
	AddRawGoFunc(env, "cons", scm_cons)
	AddRawGoFunc(env, "car", scm_car)
	AddRawGoFunc(env, "cdr", scm_cdr)
	AddRawGoFunc(env, "+", scm_add)
	AddRawGoFunc(env, "-", scm_subtract)
	AddRawGoFunc(env, "*", scm_multiplication)
	AddRawGoFunc(env, "/", scm_division)

	// Debug
	display(env)
	fmt.Println("")

	for {
		// Read
		expr := getexpr(reader)
		// Eval
		ret := eval(env, expr)
		// Print
		display(ret)
		fmt.Println("")
	}
}
