package goscm

type SCMT_Environment struct {
	table map[string]SCMT
	child *SCMT_Environment
}

func (*SCMT_Environment) scm_eval(*SCMT_Environment) SCMT {
	// TODO: This should probably be an error
	return nil
}

func (*SCMT_Environment) String() string {
	return "#<environment>"
}

func EnvEmpty(child *SCMT_Environment) *SCMT_Environment {
	return &SCMT_Environment {
		table: make(map[string]SCMT),
		child: child,
	}
}

func (env *SCMT_Environment) Add(symb *SCMT_Symbol, val SCMT) {
	env.table[symb.value] = val
}

func (env *SCMT_Environment) Find(symb *SCMT_Symbol) SCMT {
	return env.table[symb.value]
}
