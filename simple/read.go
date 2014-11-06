package simple

import (
	"github.com/JohnAnthony/goscm"
	"bufio"
	"strings"
	"strconv"
	"errors"
)

func ReadStr(str string) (goscm.SCMT, error) {
	ret, err := Read(bufio.NewReader(strings.NewReader(str)))
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	return ret, err
}

func Read(b *bufio.Reader) (goscm.SCMT, error) {
	err := chomp_meaningless(b)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	c, err := b.ReadByte()
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	switch {
	case c == '(': 
		// A list
		return read_list(b)
	case c == ')':
		// Unbalanced list
		return goscm.SCMT_Nil, errors.New("Mismatched parenthesis")
	case c == '\'':
		// A piece of quoted syntax
		subexpr, err := Read(b)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		return goscm.Make_List(
			goscm.Make_Symbol("quote"),
			subexpr), nil
	case c >= '0' && c <= '9':
		// An integer
		b.UnreadByte()
		return read_integer(b)
	case c == '"':
		// A string
		return read_string(b)
	case c == '#':
		// Something special, but for simple just a bool
		return read_bool(b)
	default:
		// A symbol
		b.UnreadByte()
		return read_symbol(b)
	}
}

func is_whitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t'
}

func read_comment(b *bufio.Reader) {
	for {
		c, err := b.ReadByte()
		switch {
		case c == '\n':
			fallthrough
		case err != nil:
			return
		default:
			continue
		}
	}
}

func read_symbol(b *bufio.Reader) (goscm.SCMT, error) {
	ret := ""
	var c byte
	var err error

	for {
		c, err = b.ReadByte()
		if err != nil {
			break
		}
		if c == ')' {
			b.UnreadByte()
			break
		}
		if c == ';' {
			read_comment(b)
			continue
		}
		if is_whitespace(c) {
			break
		}
		ret = string(append([]byte(ret), c))
	}
	return goscm.Make_Symbol(ret), err
}

func read_string(b *bufio.Reader) (goscm.SCMT, error) {
	ret := ""
	var c byte
	var err error

	for {
		c, err = b.ReadByte()
		if err != nil {
			break
		}
		if c == '"' {
			break
		}
		ret = string(append([]byte(ret), c))
	}
	return goscm.Make_SCMT(ret), err
}

func read_integer(b *bufio.Reader) (goscm.SCMT, error) {
	ret := ""
	var c byte
	var err error
	
	for {
		c, err = b.ReadByte()
		if err != nil {
			break
		}
		if c == ')' {
			b.UnreadByte()
			break
		}
		if c == ';' {
			read_comment(b)
			continue
		}
		if is_whitespace(c) {
			break
		}
		ret = string(append([]byte(ret), c))
	}
	retn, _ := strconv.Atoi(ret)
	return goscm.Make_SCMT(retn), err
}

func read_list(b *bufio.Reader) (goscm.SCMT, error) {
	c, err := b.ReadByte()
	if c == ')' {
		return goscm.SCMT_Nil, nil
	}
	if c == '.' {
		cdr, err := Read(b)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		return cdr, nil
	}

	b.UnreadByte()
	elem, err := Read(b)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	tail, err := read_list(b)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	return goscm.Cons(elem, tail), nil
}

func read_bool(b *bufio.Reader) (goscm.SCMT, error) {
	c, err := b.ReadByte()
	
	switch {
	case err != nil:
		return goscm.SCMT_Nil, err
	case c == 'f':
		return goscm.Make_SCMT(false), nil
	case c == 't':
		return goscm.Make_SCMT(true), nil
	default:
		return goscm.SCMT_Nil, errors.New("Error reading bool")
	}
}

func chomp_meaningless(b *bufio.Reader) error {
	var c byte
	var err error
	
	for {
		c, err = b.ReadByte()
		if  err != nil {
			break
		}
		if c == ';' {
			read_comment(b)
			continue
		}
		if !is_whitespace(c) {
			break
		}
	}
	b.UnreadByte()
	return err
}
