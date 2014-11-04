package goscm

import "strings"

type SCMT_Symbol struct {
	Value string
}

func (symb *SCMT_Symbol) Eval(env *SCMT_Env) (SCMT, error) {
	return env.Find(symb)
}

func (symb *SCMT_Symbol) String() string {
	return symb.Value
}

// Helpers

func Make_Symbol(s string) *SCMT_Symbol {
	return &SCMT_Symbol {
		Value: strings.ToUpper(s),
	}
}
