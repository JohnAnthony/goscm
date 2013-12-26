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

func NewInstance() *Instance {
	return &Instance {
		env: nil,
		running: true,
	}
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
	fmt.Printf("|> ")

	if expr == nil {
		fmt.Println("nil")
		return
	}
	
	switch expr.stype {
	case scm_int:
		fmt.Printf("%d :: Int", *expr.value.(*int))
	case scm_string:
		fmt.Printf("\"%s\" :: String", expr.value.(string))
	case scm_symbol:
		fmt.Printf("%s :: Symbol", expr.value.(string))
	case scm_pair:
		fmt.Printf("<#pair>")
	default:
		fmt.Printf("<#error>")
	}

	fmt.Println("")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read
		expr := getexpr(reader)
		display(expr)
		// Eval
		// Print
		
	}
}
