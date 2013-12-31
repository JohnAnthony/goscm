package goscm

///////////
// TODO: //
///////////

// Errors to the Instance through signals
// :: Consing with . mis-matched
// :: Wrong number/type of args to function
// :: Ill-formed symbols or numbers
// :: Floating-point arithmetic

// Numerical tower

///////////
// BUGS: //
///////////

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
	scm_procedure
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

func SCMString(str string) *Cell {
	return &Cell {
		stype: scm_string,
		value: str,
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

func cons(a *Cell, b *Cell) *Cell {
	return &Cell {
		stype: scm_pair,
		value: &Pair { car: a, cdr: b },
	}
}

// READ

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
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
		// TODO: This is obviously woefully inadequate and needs to handle
		// different number types
		val, _ := strconv.Atoi(str)
		return SCMInteger(val)
	}

	if str[0] == '"' {
		// TODO: The last character should be " and we should only be trimming
		// one " from each side. AND we need to handle escape characters
		return SCMString(strings.Trim(str, "\""))
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
		b, _ = r.ReadByte()
		if !isSpace(b) {
			r.UnreadByte()
			break
		}
		return tok_dot, ""
	}

	// Read our symbol into the buffer
	r.UnreadByte()
	for {
		b, _ = r.ReadByte()

		if isSpace(b) {
			break
		}

		if b == '(' || b == ')' {
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
		inst.depthAdd()
		car = inst.parse(r)
	case tok_closesub:
		inst.depthRem()
		return nil
	case tok_dot:
		tokt, tokv = gettoken(r)
		tokt, _ = gettoken(r)
		// TODO: If tokt isn't tok_closesub we've got a problem
		inst.depthRem()
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

func (inst *Instance) eval(env *Cell, expr *Cell) (nenv *Cell, ret *Cell) {
	if expr == nil {
		return env, nil
	}

	if car(expr).stype == scm_symbol {
		fmt.Printf("SYMBOL: %s\n", car(expr).value.(string))
		switch car(expr).value.(string) {
		case "quote":
			// TODO: cdr(cdr(expr)) not being nil is an error
			return env, car(cdr(expr))
		case "define":
			symb := car(cdr(expr))
			nenv, value := inst.eval(env, cdr(cdr(expr)))
			pair := cons(symb, value)
			return cons(pair, nenv), symb
		}
	}

	// Symbol lookup in environment

	//Non-special forms ...

	return env, expr
}

// PRINT

func display(c *Cell) string {
	if c == nil {
		return "nil"
	}

	switch c.stype {
	case scm_symbol:
		return c.value.(string)
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
		return fmt.Sprintf("\"%s\"", c.value.(string))
	case scm_gofunc:
		return "#<GOFUNC>"
	case scm_procedure:
		return "#<PROCEDURE>"
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
	env         *Cell
}

func NewInstance() *Instance {
	return &Instance{
		paren_depth: 0,
		env:         nil,
	}
}

func (inst *Instance) depthAdd() {
	inst.paren_depth++
}

func (inst *Instance) depthRem() {
	write := bufio.NewWriter(os.Stderr)
	inst.paren_depth--
	if inst.paren_depth < 0 {
		fmt.Fprintln(write, "ERR :: Unbalanced parentheses - too many")
		write.Flush()
	}
}

func (inst *Instance) EnvironmentalEval(expr *Cell) *Cell {
	inst.env, expr = inst.eval(inst.env, expr)
	return expr
}

func (inst *Instance) REPL(fin *os.File, fout *os.File) {
	var expr *Cell
	read := bufio.NewReader(fin)
	write := bufio.NewWriter(fout)

	for {
		////////// READ //////////
		expr = inst.parse(read)

		if inst.paren_depth > 0 {
			write := bufio.NewWriter(os.Stderr)
			fmt.Fprintln(write, "ERR :: EOF reached with unterminated parens")
			write.Flush()
			return
		}

		// Break out when parse returns nothing (EOF)
		if expr == nil {
			break
		}

		////////// EVAL //////////

		for c := expr; c != nil; c = cdr(c) {
			fmt.Fprintln(write, display(car(c)))
			expr = inst.EnvironmentalEval(car(c))

			////////// PRINT //////////

			fmt.Fprintf(write, "|> %s\n", display(expr))
			write.Flush()
		}
	}

	fmt.Printf("Ending environment: %s\n", display(inst.env))
}
