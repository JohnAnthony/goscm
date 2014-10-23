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

func EnvEmpty(parent *SCMT_Env) *SCMT_Env {
	return &SCMT_Env {
		table: make(map[string]SCMT),
		parent: parent,
	}
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
