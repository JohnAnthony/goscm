package goscm

// TODO:

// BUGS:

// Paren matching. Currently unbalanced parens seem to be wonky (A paren
// effectively just ends input like an EOF, so an early paren causes EOF) and
// EOF effecitvly just closes all parens. I actually don't know what to do about
// this

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ScmType int

const (
	scm_symbol ScmType = iota
	scm_number
	scm_complex
	scm_real
	scm_rational
	scm_integer
	scm_boolean
	scm_string
	scm_gofunc
	scm_pair
)

type Cell struct {
	stype ScmType
	value interface{}
}

type Pair struct {
	car *Cell
	cdr *Cell
}

func SCMSymbol(str string) *Cell {
	return &Cell{
		stype: scm_symbol,
		value: str,
	}
}

func SCMInteger(n int) *Cell {
	return &Cell{
		stype: scm_integer,
		value: &n,
	}
}

func SCMBoolean(b bool) *Cell {
	return &Cell{
		stype: scm_boolean,
		value: &b,
	}
}

func SCMPair(a *Cell, b *Cell) *Cell {
	return &Cell{
		stype: scm_pair,
		value: &Pair{car: a, cdr: b},
	}
}

// Scheme-like functions

func car(c *Cell) *Cell {
	return c.value.(*Pair).car
}

func cdr(c *Cell) *Cell {
	return c.value.(*Pair).cdr
}

// READ

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	}
	return false
}

func isSpecial(b byte) bool {
	switch b {
	case '(', ')':
		return true
	}
	return false
}

type ScmTok int

const (
	tok_identifier ScmTok = iota
	tok_opensub
	tok_closesub
	tok_eof
	tok_dot
)

func identifier_to_cell(str string) *Cell {
	if str[0] >= '0' && str[0] <= '9' {
		val, _ := strconv.Atoi(str)
		return SCMInteger(val)
	}

	if strings.ToLower(str) == "#t" {
		return SCMBoolean(true)
	}

	if strings.ToLower(str) == "#f" {
		return SCMBoolean(false)
	}

	return SCMSymbol(str)
}

func gettoken(r *bufio.Reader) (t ScmTok, value string) {
	var b byte
	var err error
	buffer := make([]byte, 0)

	// Chomp whitespace
	for {
		b, err = r.ReadByte()
		if err == io.EOF {
			return tok_eof, ""
		}

		if !isSpace(b) {
			break
		}
	}

	switch b {
	case '(':
		return tok_opensub, ""
	case ')':
		return tok_closesub, ""
	case '.':
		// TODO: Need to check for a space afterwards
		return tok_dot, ""
	}

	// Read our symbol into the buffer
	r.UnreadByte()
	for {
		b, _ = r.ReadByte()

		if isSpace(b) {
			break
		}

		if isSpecial(b) {
			r.UnreadByte()
			break
		}
		buffer = append(buffer, b)
	}

	return tok_identifier, string(buffer)
}

func (inst *Instance) parse(r *bufio.Reader) *Cell {
	var car *Cell
	tokt, tokv := gettoken(r)
	switch tokt {
	case tok_identifier:
		car = identifier_to_cell(tokv)
	case tok_opensub:
		inst.paren_depth++
		car = inst.parse(r)
	case tok_closesub:
		inst.paren_depth--
		return nil
	case tok_dot:
		tokt, tokv = gettoken(r)
		tokt, _ = gettoken(r)
		// TODO: If tokt isn't tok_closesub we've got a problem
		return identifier_to_cell(tokv)
	case tok_eof:
		return nil
	}

	// If we are in the bottom-level interpreter
	if inst.paren_depth == 0 {
		return SCMPair(car, nil)
	}

	return SCMPair(car, inst.parse(r))
}

// EVAL

// PRINT

func display(c *Cell) string {
	if c == nil {
		return "nil"
	}

	switch c.stype {
	case scm_symbol:
		return "#<SYMBOL " + c.value.(string) + ">"
	case scm_number:
		return "#<NUMBER>"
	case scm_complex:
		return "#<COMPLEX>"
	case scm_real:
		return "#<REAL>"
	case scm_rational:
		return "#<RATIONAL>"
	case scm_integer:
		return fmt.Sprintf("%d", *c.value.(*int))
	case scm_boolean:
		if *c.value.(*bool) == true {
			return "#t"
		}
		return "#f"
	case scm_string:
		return "#<STRING>"
	case scm_gofunc:
		return "#<GOFUNC>"
	case scm_pair:
		str := "("
		for ; c != nil; c = cdr(c) {
			str += display(car(c))
			if cdr(c) == nil {
				break
			}
			if cdr(c).stype != scm_pair {
				str += " . " + display(cdr(c))
				break
			}
			if cdr(c) != nil {
				str += " "
			}
		}
		str += ")"
		return str
	}

	// We should never be getting here
	return "#<ERROR>"
}

// EXPORTED

type Instance struct {
	paren_depth int
	env *Cell
}

func NewInstance() *Instance {
	return &Instance {
		paren_depth: 0,
		env: nil,
	}
}

func (inst* Instance) REPL(fin *os.File, fout *os.File) {
	var expr *Cell
	read := bufio.NewReader(fin)
	write := bufio.NewWriter(fout)

	// READ
	for {
		expr = inst.parse(read)
		if expr == nil {
			break
		}

		if inst.paren_depth != 0 {
			fmt.Fprintln(write, "ERR: Unbalanced parentheses")
			write.Flush()
			return
		}

		// EVAL
		//		expr = eval(expr)
		
		// PRINT
		for c := expr; c != nil; c = cdr(c) {
			fmt.Fprintf(write, "|> %s\n", display(car(c)))
		}
		write.Flush()
	}
}
