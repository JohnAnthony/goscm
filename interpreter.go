package main

// TODO:
// Comments
// quote '
// quasiquote `
// unquote ,
// scm_float
// Escaped characters
// Tail Call Optimisation
// Numerical tower (number / complex / real / rational / integer
// cons infix notation .
// SPECIAL FORMS:
// let and let*
// SCM FUNCTIONS:
// eq? eqv? equal?
// atom? list?
// zero? false? true? not
// or and
// apply
// begin
// string=?

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ScmType int

type ScmPair struct {
	car *Cell
	cdr *Cell
}

type Cell struct {
	stype ScmType
	value interface{}
}

const (
	scm_bool ScmType = iota
	scm_emptylist
	scm_func
	scm_gofunc
	scm_int
	scm_pair
	scm_specialform
	scm_string
	scm_symbol
)

func SCMBool(b bool) *Cell {
	return &Cell {
		stype: scm_bool,
		value: &b,
	}
}

func SCMEmptyList() *Cell {
	return &Cell {
		stype: scm_emptylist,
	}
}

func SCMFunc(name string, f *Cell) *Cell {
	return &Cell {
		stype: scm_func,
		value: SCMPair(SCMSymbol(name), f),
	}
}

func SCMGoFunc(name string, f func (*Cell) *Cell) *Cell {
	fcell := Cell {
		stype: scm_gofunc,
		value: f,
	}
	return SCMPair(SCMSymbol(name), &fcell)
}

func SCMInt(n int) *Cell {
	return &Cell {
		stype: scm_int,
		value: &n,
	}
}

func SCMPair(a *Cell, b *Cell) *Cell {
	return &Cell{
		stype: scm_pair,
		value: &ScmPair{car: a, cdr: b},
	}
}

func SCMSpecialForm(name string, f func (*Cell) *Cell) *Cell {
	fcell := Cell {
		stype: scm_specialform,
		value: f,
	}
	return SCMPair(SCMSymbol(name), &fcell)
}

func SCMString(str string) *Cell {
	return &Cell {
		stype: scm_string,
		value: str,
	}
}
func SCMSymbol(str string) *Cell {
	return &Cell {
		stype: scm_symbol,
		value: strings.ToUpper(str),
	}
}

func NewEnvironment() *Cell {
	var env *Cell = SCMEmptyList()
	// Special forms
	env = addSpecialForm(env, "quote", scm_quote)
	env = addSpecialForm(env, "define", scm_define)
	env = addSpecialForm(env, "lambda", scm_lambda)
	env = addSpecialForm(env, "if", nil)
	env = addSpecialForm(env, "load-from-path", nil)
	// Go functions
	env = AddRawGoFunc(env, "cons", scm_cons)
	env = AddRawGoFunc(env, "car", scm_car)
	env = AddRawGoFunc(env, "cdr", scm_cdr)
	env = AddRawGoFunc(env, "display", scm_display)
	env = AddRawGoFunc(env, "+", scm_add)
	env = AddRawGoFunc(env, "-", scm_subtract)
	env = AddRawGoFunc(env, "*", scm_multiplication)
	env = AddRawGoFunc(env, "/", scm_division)
	env = AddRawGoFunc(env, "=", scm_numeq)
	return env
}

// recursively duplicate scheme objects
func duplicate(cell *Cell) *Cell {
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
	case scm_specialform:
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
		return SCMPair(duplicate(car(cell)), duplicate(cdr(cell)))
	}

	// If we get here it's an error
	return nil
}

func eval(env *Cell, expr *Cell) (newenv *Cell, result *Cell) {
	switch expr.stype {
	case scm_int:
		fallthrough
	case scm_string:
		fallthrough
	case scm_gofunc:
		return env, expr
	case scm_symbol:
		return env, symbolLookup(env, expr.value.(string))
	case scm_func:
		return env, expr
	case scm_emptylist:
		return env, expr
	case scm_bool:
		return env, expr
	case scm_specialform:
		return env, expr
	case scm_pair:
		// otherwise, we're using golang functions
		funcsym := car(expr).value.(string)
		f := symbolLookup(env, funcsym)
		if f == nil {   // Symbol not found
			return env, nil
		}
		tail := duplicate(cdr(expr))

		// If special case
		if f.stype == scm_specialform && funcsym == "IF" {
			env, ev := eval(env, car(tail))
			pred := ev.value.(*bool)
			fst := car(cdr(tail))
			snd := car(cdr(cdr(tail)))
			var result *Cell
			
			if *pred == true {
				_, result = eval(env, fst)
			} else {
				_, result = eval(env, snd)
			}
			return env, result
		}

		// load-from-path special form
		if f.stype == scm_specialform && funcsym == "LOAD-FROM-PATH" {
			path := car(tail).value.(string)
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return env, SCMBool(false)
			}
			defer file.Close()
			return ReadFile(env, file)
		}

		// We don't eval if quoting
		if f.stype != scm_specialform || funcsym == "DEFINE" {
			var e *Cell

			// Special case for define - we don't eval the first symbol
			if f.stype == scm_specialform && funcsym == "DEFINE" {
				e = cdr(tail)
			} else {
				e = tail
			}

			// Not a special case - eval everything and build a new "tail"
			for ; e.stype != scm_emptylist; e = cdr(e) {
				_, e.value.(*ScmPair).car = eval(env, car(e))
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
				_, ret = eval(subenv, car(sub))
			}

			return env, ret
		}

		// If we reach this point we should be handling scm_gofunc and scm_specialform types

		ret := f.value.(func(*Cell) *Cell)(tail)

		// Handle define
		if funcsym == "DEFINE" {
			// Shuffle our new define into the place of the old env head
			tmpenv := *env
			*env = *cons(ret, &tmpenv)

			return env, SCMBool(true)
		}

		// Function is not define, return normally
		return env, ret
	}

	// Getting here is an error
	return env, nil
}

func addSpecialForm(env *Cell, symb string, f func(*Cell) *Cell) *Cell {
	newcell := SCMSpecialForm(symb, f)
	return cons(newcell, env)
}

func AddRawGoFunc(env *Cell, symb string, f func(*Cell) *Cell) *Cell {
	newcell := SCMGoFunc(symb, f)
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
	if c == 0 { // EOF
		return nil
	}

	// Sytax keyed to the first character
	switch c {
	case '\n':        // After chomping a '\n' means we're sticking on EOF
		return nil
	case ')':      // End of a list
		return SCMEmptyList()
	case '(':         // A list
		var head *Cell = nil
		nexp := getexpr(in)
		if nexp.stype == scm_emptylist {
			return SCMEmptyList()
		}
		
		head = &Cell{
			stype: scm_pair,
			value: &ScmPair{car: nexp},
		}
		tip := head

		for {
			nexp = getexpr(in)
			if nexp.stype == scm_emptylist {
				tip.value.(*ScmPair).cdr = SCMEmptyList()
				break
			}
			tip.value.(*ScmPair).cdr = &Cell{
				stype: scm_pair,
				value: &ScmPair{car: nexp},
			}
			tip = cdr(tip)
		}
		return head
	case '"':      // A string
		for c, _ = in.ReadByte(); c != '"'; c, _ = in.ReadByte() {
			symbol = append(symbol, c)
		}
		return SCMString(string(symbol))
	}

	// Alright, read the whole symbol
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
		return SCMInt(n)
	}

	// Symbol
	return SCMSymbol(string(symbol))
}

func display(expr *Cell) {
	if expr == nil {
//		fmt.Printf("nil")
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
	case scm_specialform:
		fmt.Printf("#<specialform>")
	default:
		fmt.Printf("#<error>")
	}
}

// Go versions of scm functions

func cons(a *Cell, b *Cell) *Cell {
	return SCMPair(a, b)
}

func car(lst *Cell) *Cell {
	return lst.value.(*ScmPair).car
}

func cdr(lst *Cell) *Cell {
	return lst.value.(*ScmPair).cdr
}

// Scm functions
// Special forms
func scm_quote(tail *Cell) *Cell {
	return car(tail)
}

func scm_define(tail *Cell) *Cell {
	return SCMPair(car(tail), car(cdr(tail)))
}

func scm_lambda(tail *Cell) *Cell {
    return &Cell{
	    stype: scm_func,
		value: &ScmPair{car: car(tail), cdr: cdr(tail)},
	}
}

// Non-special form functions
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

	return SCMInt(ret)
}

func scm_subtract(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e.stype != scm_emptylist; e = cdr(e) {
		ret -= *car(e).value.(*int)
	}

	return SCMInt(ret)
}

func scm_multiplication(tail *Cell) *Cell {
	ret := 1

	for e := tail; e.stype != scm_emptylist; e = cdr(e) {
		ret *= *car(e).value.(*int)
	}

	return SCMInt(ret)
}

func scm_division(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e.stype != scm_emptylist; e = cdr(e) {
		ret /= *car(e).value.(*int)
	}

	return SCMInt(ret)
}

func scm_numeq(tail *Cell) *Cell {
	fst := *car(tail).value.(*int)
	
	for e := cdr(tail); e.stype != scm_emptylist; e = cdr(e) {
		if *car(e).value.(*int) != fst {
			return SCMBool(false)
		}
	}

	return SCMBool(true)
}

// External interface

func ReadFile(env *Cell, file *os.File) (newenv *Cell, result *Cell) {
	reader := bufio.NewReader(file)
	subenv := env
	for {
		// Read
		expr := getexpr(reader)
		if expr == nil {
			break
		}
		// Eval
		subenv2, ret := eval(subenv, expr)
		subenv = subenv2
		// Print
		display(ret)
		fmt.Println("")
	}
	return subenv, SCMBool(true)
}

func main() {
	env := NewEnvironment()
	ReadFile(env, os.Stdin)
}
