package simple

import (
	"github.com/JohnAnthony/goscm"
	"io"
	"bufio"
)

func Read(r io.Reader) goscm.SCMT {
	b := bufio.NewReader(r)

	// Chomp preceeding whitespace
	for {
		c, err := b.ReadByte()
		if err != nil {
			break
		}
		if c != ' ' {
			break
		}
	}
	b.UnreadByte()
	
	switch {

	default: // A symbol
		ret := ""
		for {
			c, err := b.ReadByte()
			if err != nil {
				break
			}
			if c == ' ' {
				break
			}
			ret = string(append([]byte(ret), c))
		}
		return goscm.Make_Symbol(ret)
	}
}
