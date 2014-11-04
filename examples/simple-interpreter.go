package main

import (
	"github.com/JohnAnthony/goscm"
	"github.com/JohnAnthony/goscm/simple"
	"bufio"
	"os"
)

func main() {
	goscm.REPL(bufio.NewReader(os.Stdin), simple.Read, simple.Env())
}
