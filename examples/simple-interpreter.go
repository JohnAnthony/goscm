package main

import (
	"github.com/JohnAnthony/goscm"
	"github.com/JohnAnthony/goscm/simple"
	"bufio"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	goscm.REPL(in, simple.Read, simple.Env())
}
