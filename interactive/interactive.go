package main

import (
	"goscm"
	"os"
)

func main() {
	inst := goscm.NewInstance()
	inst.REPL(os.Stdin, os.Stdout)
}
