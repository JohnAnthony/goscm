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

func toktype(tok string) ScmType {
	if tok[0] == '"' {
		return scm_string
	}
	if tok[0] >= '0' && tok[0] <= '9' {
		return scm_int
	}
	return scm_symbol
}

func gettoken(in *bufio.Reader) string {
	var c byte
	symbol := make([]byte, 0)
	
	// Skip over whitespace
	for {
		c, _ = in.ReadByte()
		if !isSpace(c) {
			break
		}
	}

	if c == '(' {
		return "("
	}

	if c == ')' {
		return ")"
	}

	for !isSpace(c) && c != ')' {
		symbol = append(symbol, c)
		c, _ = in.ReadByte()
	}

	if c == ')' {
		in.UnreadByte()
	}

	return string(symbol)
}

func getexpression(r *bufio.Reader) *Cell {
	tok := gettoken(r)
	t := toktype(tok)
	cell := Cell {
		stype: t,
	}

	switch t {
	case scm_int:
		n, err := strconv.Atoi(tok)
		if err != nil {
			fmt.Println(err)
		}
		cell.value = &n
	case scm_string:
		cell.value = tok
	case scm_symbol:
		cell.value = tok
	}

	return &cell
}

func display(expr *Cell) {
	fmt.Printf("|> ")

	switch expr.stype {
	case scm_int:
		fmt.Printf("%d :: Int\n", *expr.value.(*int))
	case scm_string:
		fmt.Printf("\"%s\" :: String\n", expr.value.(string))
	case scm_symbol:
		fmt.Printf("%s :: Symbol\n", expr.value.(string))
	default:
		fmt.Println("<#error>")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read
		expr := getexpression(reader)
		display(expr)
		// Eval
		// Print
		
	}
}
