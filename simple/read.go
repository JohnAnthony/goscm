package simple

import (
	"github.com/JohnAnthony/goscm"
	"io"
	"bufio"
	"strconv"
)

func Read(r io.Reader) goscm.SCMT {
	b := bufio.NewReader(r)
	var c byte
	var err error

	// Chomp preceeding whitespace
	for {
		c, err = b.ReadByte()
		if err != nil {
			break
		}
		if c != ' ' {
			break
		}
	}
	
	switch {
	case c >= '0' && c <= '9': // An integer
		ret := ""
		for {
			if err != nil {
				break
			}
			if c == ' ' {
				break
			}
			ret = string(append([]byte(ret), c))
			c, err = b.ReadByte()
		}
		retn, _ := strconv.Atoi(ret)
		return goscm.Make_SCMT(retn)
	default: // A symbol
		ret := ""
		for {
			if err != nil {
				break
			}
			if c == ' ' {
				break
			}
			ret = string(append([]byte(ret), c))
			c, err = b.ReadByte()
		}
		return goscm.Make_Symbol(ret)
	}
}
