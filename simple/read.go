package simple

import (
	"github.com/JohnAnthony/goscm"
	"bufio"
	"strings"
	"strconv"
)

func ReadStr(str string) (goscm.SCMT, error) {
	ret, err := Read(bufio.NewReader(strings.NewReader(str)))
	return ret, err
}

func Read(b *bufio.Reader) (goscm.SCMT, error) {
	chomp_meaningless(b)
	c, err := b.ReadByte()
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	switch {
	case c == '(': 
		// A list
		return read_list(b)
	case c >= '0' && c <= '9':
		// An integer
		b.UnreadByte()
		return read_integer(b)
	case c == '"':
		// A string
		return read_string(b)
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
	list := goscm.SCMT_Nil
	var c byte
	var err error

	for {
		c, err = b.ReadByte()
		if c == ')' {
			break
		}
		b.UnreadByte()
		recurse, _ := Read(b)
		list = goscm.Cons(recurse, list)
	}
	return goscm.Reverse(list), err
}

func chomp_meaningless(b *bufio.Reader) {
	for {
		c, err := b.ReadByte()
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
}
