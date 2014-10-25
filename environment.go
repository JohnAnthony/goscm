package goscm

type SCMT_Env struct {
	table map[string]SCMT
	parent *SCMT_Env
}

func (env *SCMT_Env) scm_eval(*SCMT_Env) SCMT {
	return env
}

func (*SCMT_Env) String() string {
	return "#<environment>"
}

func (env *SCMT_Env) Add(symb *SCMT_Symbol, val SCMT) {
	env.table[symb.value] = val
}

func (env *SCMT_Env) Find(symb *SCMT_Symbol) SCMT {
	ret := env.table[symb.value]
	if ret != nil {
		return ret
	}
	if env.parent == nil {
		return nil
	}
	return env.parent.Find(symb)
}

// Environment provided #1: Completely empty. This is used for procedure
// environment closures.
// Procedures provided: NONE

func EnvEmpty(parent *SCMT_Env) *SCMT_Env {
	return &SCMT_Env {
		table: make(map[string]SCMT),
		parent: parent,
	}
}

// Environment provided #2: Simple. This is helpful for testing and as an
// example of how to build your own scheme env from scratch.
// Procedures provided: + - / * quote car cdr cons

func EnvSimple() *SCMT_Env {
	env := EnvEmpty(nil)
//	env.Add(Make_Symbol("+", scm_add))
//	env.Add(Make_Symbol("-", scm_subtract)
//	env.Add(Make_Symbol("/", scm_divide))
//	env.Add(Make_Symbol("*", scm_multiply))
//	env.Add(Make_Symbol("quote", scm_quote))
//	env.Add(Make_Symbol("car", scm_car))
//	env.Add(Make_Symbol("cdr", scm_cdr))
//	env.Add(Make_Symbol("cons", scm_cons))
	return env
}
