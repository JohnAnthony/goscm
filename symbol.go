package goscm

import "strings"

type SCMT_Symbol struct {
	value string
}

func (symb *SCMT_Symbol) scm_eval(env *SCMT_Env) SCMT {
	// TODO: Look up in environment!
	return nil
}

func (symb *SCMT_Symbol) String() string {
	return symb.value
}

// Helpers

func Make_Symbol(s string) *SCMT_Symbol {
	return &SCMT_Symbol {
		value: strings.ToUpper(s),
	}
}
