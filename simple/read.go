package simple

import (
	"github.com/JohnAnthony/goscm"
	"bufio"
	"strings"
	"strconv"
)

func ReadStr(str string) goscm.SCMT {
	ret, _ := Read(bufio.NewReader(strings.NewReader(str)))
	return ret
}

func Read(b *bufio.Reader) (goscm.SCMT, *bufio.Reader) {
	var c byte
	var err error

	// Chomp preceeding whitespace
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
	
	switch {
	case c == '(': // A list
		list := goscm.SCMT_Nil
		for {
			c, err = b.ReadByte()
			if c == ')' {
				break
			}
			b.UnreadByte()
			recurse, _ := Read(b)
			list = goscm.Cons(recurse, list)
		}
		return goscm.Reverse(list), b
	case c >= '0' && c <= '9': // An integer
		ret := ""
		for {
			if err != nil {
				break
			}
			if c == ')' {
				b.UnreadByte()
				break
			}
			if is_whitespace(c) {
				break
			}
			ret = string(append([]byte(ret), c))
			c, err = b.ReadByte()
		}
		retn, _ := strconv.Atoi(ret)
		return goscm.Make_SCMT(retn), b
	case c == '"': // A string
		ret := ""
		for {
			c, err = b.ReadByte()
			if c == '"' {
				break
			}
			ret = string(append([]byte(ret), c))
		}
		return goscm.Make_SCMT(ret), b
	default: // A symbol
		ret := ""
		for {
			if err != nil {
				break
			}
			if c == ')' {
				b.UnreadByte()
				break
			}
			if is_whitespace(c) {
				break
			}
			ret = string(append([]byte(ret), c))
			c, err = b.ReadByte()
		}
		return goscm.Make_Symbol(ret), b
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
