package main

import (
	"os"
	"bufio"
	"fmt"
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

type Instance struct {
	env *Cell
	running bool
}

func (inst *Instance) eval(expr *Cell) *Cell {
	if expr.stype != scm_pair {
		return expr
	}

	funcsym := car(expr).value.(string)
	tail := cdr(expr)
	f := symbolLookup(inst.env, funcsym)

	if f == nil {
		return nil
	}

	// We don't eval if quoting!
	if funcsym != "quote" {
		// Eval all the tail args first
		for e := tail; e != nil; e = cdr(e) {
			e.value.(*ScmPair).car = inst.eval(car(e))
		}
	}
		
	// Go source functions
	if f.stype == scm_gofunc {
		return f.value.(func (*Cell) *Cell)(tail)
	}

	// ??? if scm_func type

	// Error because we've got an unhandled type
	return nil
}

func NewInstance() *Instance {
	return &Instance {
		env: nil,
		running: true,
	}
}

func (inst *Instance) AddRawGoFunc(symb string, f func(*Cell) *Cell) {
	newcar := &Cell {
		stype: scm_symbol,
		value: symb,
	}
	newcdr := &Cell {
		stype: scm_gofunc,
		value: f,
	}
	newcell := cons(newcar, newcdr)
	inst.env = cons(newcell, inst.env)
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

		head = &Cell {
			stype: scm_pair,
			value: &ScmPair { car: nexp },
		}
		tip := head
		for {
			nexp = getexpr(in)
			if nexp.stype == scm_pairclose {
				tip.value.(*ScmPair).cdr = nil
				break
			}
			tip.value.(*ScmPair).cdr = &Cell {
				stype: scm_pair,
				value: &ScmPair { car: nexp },
			}
			tip = cdr(tip)
		}
		return head
	}

	// Scm Pair Close
	if c == ')' {
		return &Cell {
			stype: scm_pairclose,
			value: ")",
		}
	}

	// Handle strings
	if c == '"' {
		for c, _ = in.ReadByte(); c != '"'; c, _ = in.ReadByte() {
			symbol = append(symbol, c)
		}
		return &Cell {
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
		
		return &Cell {
			stype: scm_int,
			value: &n,
		}
	}
	
	// Symbol
	return &Cell {
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
		// This doesn't actually display pairs properly
		// TODO: Fix that
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
	return &Cell {
		stype: scm_pair,
		value: &ScmPair { car: a, cdr: b },
	}
}

func car(lst *Cell) *Cell {
	return lst.value.(*ScmPair).car
}

func cdr(lst *Cell) *Cell {
	return lst.value.(*ScmPair).cdr
}

// Scm functions

// scm_quit function
// scm_define
// scm_lambda

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

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_subtract(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e != nil; e = cdr(e) {
		ret -= *car(e).value.(*int)
	}

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_multiplication(tail *Cell) *Cell {
	ret := 1

	for e := tail; e != nil; e = cdr(e) {
		ret *= *car(e).value.(*int)
	}

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_division(tail *Cell) *Cell {
	ret := *car(tail).value.(*int)

	for e := cdr(tail); e != nil; e = cdr(e) {
		ret /= *car(e).value.(*int)
	}

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_quote(tail *Cell) *Cell {
	// TODO: Syntax checking. Quote takes exactly one argument
	return car(tail)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	inst := NewInstance()
	inst.AddRawGoFunc("quote", scm_quote)
	inst.AddRawGoFunc("cons", scm_cons)
	inst.AddRawGoFunc("car", scm_car)
	inst.AddRawGoFunc("cdr", scm_cdr)
	inst.AddRawGoFunc("+", scm_add)
	inst.AddRawGoFunc("-", scm_subtract)
	inst.AddRawGoFunc("*", scm_multiplication)
	inst.AddRawGoFunc("/", scm_division)

	for inst.running {
		// Read
		expr := getexpr(reader)
		// Eval
		ret := inst.eval(expr)
		// Print
		display(ret)
		fmt.Println("")
	}
}
