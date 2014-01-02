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
		value: strings.ToLower(str),
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
	return &Cell{
		stype: scm_string,
		value: str,
	}
}

func SCMGoFunc(f func(*Cell) *Cell) *Cell {
	return &Cell {
		stype: scm_gofunc,
		value: f,
	}
}

func SCMProcedure(params *Cell, body *Cell) *Cell {
	return &Cell{
		stype: scm_procedure,
		value: &Pair{car: params, cdr: body},
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
	return &Cell{
		stype: scm_pair,
		value: &Pair{car: a, cdr: b},
	}
}

func reverse(a *Cell) *Cell {
	var ret *Cell = nil
	for e := a; e != nil; e = cdr(e) {
		ret = cons(car(e), ret)
	}
	return ret
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
	tok_dot
	tok_quote
	tok_quasiquote
	tok_unquote
	tok_eof
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
	case '\'':
		return tok_quote, "'"
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
	var this *Cell
	tokt, tokv := gettoken(r)
	switch tokt {
	case tok_identifier:
		this = identifier_to_cell(tokv)
	case tok_opensub:
		inst.depthAdd()
		this = inst.parse(r)
	case tok_closesub:
		inst.depthRem()
		return nil
	case tok_dot:
		tokt, tokv = gettoken(r)
		tokt, _ = gettoken(r)
		// TODO: If tokt isn't tok_closesub we've got a problem
		inst.depthRem()
		return identifier_to_cell(tokv)
	case tok_quote:
		// This feels off
		rest := inst.parse(r)
		carlst := SCMPair(SCMSymbol("quote"), SCMPair(car(rest), nil))
		return SCMPair(carlst, cdr(rest))
		// TODO
	case tok_quasiquote:
		// TODO
	case tok_unquote:
		// TODO
	case tok_eof:
		return nil
	}

	// If we are in the bottom-level interpreter return now
	if inst.paren_depth == 0 {
		return SCMPair(this, nil)
	}

	return SCMPair(this, inst.parse(r))
}

// EVAL

func symbolLookup(env *Cell, symb string) *Cell {
	for c := env; c != nil; c = cdr(c) {
		if car(car(c)).value.(string) == symb {
			return cdr(car(c))
		}
	}
	return nil
}

func (inst *Instance) eval(env *Cell, expr *Cell) (nenv *Cell, ret *Cell) {
	if expr == nil {
		return env, nil
	}

	switch expr.stype {
	case scm_symbol:
		return env, symbolLookup(env, expr.value.(string))
	case scm_number:
		fallthrough
	case scm_complex:
		fallthrough
	case scm_real:
		fallthrough
	case scm_rational:
		fallthrough
	case scm_integer:
		fallthrough
	case scm_boolean:
		fallthrough
	case scm_string:
		fallthrough
	case scm_gofunc:
		fallthrough
	case scm_procedure:
		return env, expr
	case scm_pair:
		// Nothing
	}

	// We ONLY deal with evaluating lists from this point onwards

	head := car(expr)
	tail := cdr(expr)

	// Special form symbols
	// TODO: Move as many of the below as possible into discrete functions
	if head.stype == scm_symbol {
		switch head.value.(string) {
		case "quote":
			// TODO: Check exactly one argument
			return env, car(tail)
		case "define":
			// TODO: Check exactly two arguments
			// TODO: Type checking
			symb := car(tail)
			_, value := inst.eval(env, car(cdr(tail)))
			pair := cons(symb, value)
			return cons(pair, env), symb
		case "lambda":
			// TODO: Type checking
			return env, SCMProcedure(car(tail), cdr(tail))
		case "set!":
			// TODO: Check exactly two arguments
			// TODO: Type checking
			symb := symbolLookup(env, car(tail).value.(string))
			_, ret = inst.eval(env, car(cdr(tail)))
			*symb = *ret
			return env, nil
		default:
			env, head = inst.eval(env, head)
		}
	}

	if head == nil {
		fmt.Println("Error: Symbol not found")
		return env, nil
	}

	if head.stype == scm_procedure {
		return inst.apply(env, head, tail)
	} else if head.stype == scm_gofunc {
		return inst.goapply(env, head, tail)
	}

	// Error
	fmt.Println("Error: reached end of eval")
	return env, nil
}

func (inst *Instance) apply(env *Cell, head *Cell, tail *Cell) (nenv *Cell, ret *Cell) {
	for k, v := car(head), tail; k != nil && v != nil; k, v = cdr(k), cdr(v) {
		env = cons(SCMPair(car(k), car(v)), env)
	}
	nenv = env
	for subex := cdr(head); subex != nil; subex = cdr(subex) {
		nenv, ret = inst.eval(nenv, car(subex))
	}
	return env, ret
}

func (inst *Instance) goapply(env *Cell, head *Cell, tail *Cell) (nenv *Cell, ret *Cell) {
	f := head.value.(func (*Cell) *Cell)
	nenv = env
	var collect *Cell
	for e := tail; e != nil; e = cdr(e) {
		nenv, collect = inst.eval(nenv, car(e))
		ret = cons(collect, ret)
	}
	ret = reverse(ret)
	return env, f(ret)
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

// BUILTIN

func scm_add(tail *Cell) *Cell {
	// TODO: Check for type correctness
	value := 0
	for e := tail; e != nil; e = cdr(e) {
		value += *car(e).value.(*int)
	}
	return SCMInteger(value)
}

func scm_multiply(tail *Cell) *Cell {
	// TODO: Check for type correctness
	value := 1
	for e := tail; e != nil; e = cdr(e) {
		value *= *car(e).value.(*int)
	}
	return SCMInteger(value)
}

func scm_subtract(tail *Cell) *Cell {
	// TODO: Check for type correctness
	// TODO: Check number of arguments at least one
	value := *car(tail).value.(*int)
	for e := cdr(tail); e != nil; e = cdr(e) {
		value -= *car(e).value.(*int)
	}
	return SCMInteger(value)
}

func scm_divide(tail *Cell) *Cell {
	// TODO: Check for type correctness
	// TODO: Check number of arguments at least one
	value := *car(tail).value.(*int)
	for e := cdr(tail); e != nil; e = cdr(e) {
		value /= *car(e).value.(*int)
	}
	return SCMInteger(value)
}

func scm_car(tail *Cell) *Cell {
	// TODO: Check type correctness
	// TODO: Check number of arguments exactly one
	return car(car(tail))
}

func scm_cdr(tail *Cell) *Cell {
	// TODO: Check type correctness
	// TODO: Check number of arguments exactly one
	return cdr(car(tail))
}

func scm_begin(tail *Cell) *Cell {
	var e *Cell
	for e = tail; cdr(e) != nil; e = cdr(e) {}
	return car(e)
}

// EXPORTED

type Instance struct {
	paren_depth int
	env         *Cell
}

func NewInstance() *Instance {
	inst := Instance{
		paren_depth: 0,
		env:         nil,
	}
	inst.AddRawGoFunc("+", scm_add)
	inst.AddRawGoFunc("*", scm_multiply)
	inst.AddRawGoFunc("-", scm_subtract)
	inst.AddRawGoFunc("/", scm_divide)
	inst.AddRawGoFunc("car", scm_car)
	inst.AddRawGoFunc("cdr", scm_cdr)
	inst.AddRawGoFunc("begin", scm_begin)
	return &inst
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

func (inst *Instance) AddRawGoFunc(name string, f func(*Cell) *Cell) {
	inst.env = cons(SCMPair(SCMSymbol(name), SCMGoFunc(f)), inst.env)
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
			expr = inst.EnvironmentalEval(car(c))
			////////// PRINT //////////
			fmt.Fprintf(write, "|> %s\n", display(expr))
			write.Flush()
		}
	}

	fmt.Printf("Ending environment: %s\n", display(inst.env))
}
