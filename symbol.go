package goscm

import "strings"

type Symbol struct {
	Value string
}

func (symb *Symbol) Eval(env *Environ) (SCMT, error) {
	return env.Find(symb)
}

func (symb *Symbol) String() string {
	return symb.Value
}

// Helpers

func NewSymbol(s string) *Symbol {
	return &Symbol {
		Value: strings.ToUpper(s),
	}
}
