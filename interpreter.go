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

	funcsym := expr.value.(*ScmPair).car
	tail := expr.value.(*ScmPair).cdr
	f := symbolLookup(inst.env, funcsym.value.(string))

	if f == nil {
		return nil
	}

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
			tip = tip.value.(*ScmPair).cdr
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
		for e := expr; e != nil; e = e.value.(*ScmPair).cdr {
			display(e.value.(*ScmPair).car)
			if e.value.(*ScmPair).cdr != nil {
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

// Scm functions

// scm_cons function
// scm_car function
// scm_cdr function
// scm_quit function
// scm_define
// scm_lambda

func scm_add(tail *Cell) *Cell {
	ret := 0

	for e := tail; e != nil; e = e.value.(*ScmPair).cdr {
		ret += *e.value.(*ScmPair).car.value.(*int)
	}

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_subtract(tail *Cell) *Cell {
	ret := *tail.value.(*ScmPair).car.value.(*int)

	for e := tail.value.(*ScmPair).cdr; e != nil; e = e.value.(*ScmPair).cdr {
		ret -= *e.value.(*ScmPair).car.value.(*int)
	}

	return &Cell {
		stype: scm_int,
		value: &ret,
	}
}

func scm_quote(tail *Cell) *Cell {
	// TODO: Syntax checking. Quote takes exactly one argument
	return tail.value.(*ScmPair).car
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	inst := NewInstance()
	inst.AddRawGoFunc("quote", scm_quote)
	inst.AddRawGoFunc("+", scm_add)
	inst.AddRawGoFunc("-", scm_subtract)

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
