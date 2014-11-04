package goscm

import "errors"

type SCMT_Env struct {
	table map[string]SCMT
	parent *SCMT_Env
}

func (env *SCMT_Env) Eval(*SCMT_Env) (SCMT, error) {
	return env, nil
}

func (*SCMT_Env) String() string {
	return "#<environment>"
}

func (env *SCMT_Env) Add(symb *SCMT_Symbol, val SCMT) {
	env.table[symb.Value] = val
}

func (env *SCMT_Env) Find(symb *SCMT_Symbol) (SCMT, error) {
	ret := env.table[symb.Value]
	if ret != nil {
		return ret, nil
	}
	if env.parent == nil {
		return nil, errors.New("Unbound variable: " + symb.String())
	}
	return env.parent.Find(symb)
}

func (env *SCMT_Env) Set(symb *SCMT_Symbol, val SCMT) error {
	ret := env.table[symb.Value]
	if ret != nil {
		env.table[symb.Value] = val
		return nil
	}
	if env.parent == nil {
		return errors.New("Unbound variable: " + symb.String())
	}
	return env.parent.Set(symb, val)
}

func (env *SCMT_Env) BindForeign(name string, f func (*SCMT_Pair, *SCMT_Env) (SCMT, error)) {
	env.Add(Make_Symbol(name), Make_Foreign(f))
}

func (env *SCMT_Env) BindSpecial(name string, f func (*SCMT_Pair, *SCMT_Env) (SCMT, error)) {
	env.Add(Make_Symbol(name), Make_Special(f))
}

func EnvEmpty(parent *SCMT_Env) *SCMT_Env {
	return &SCMT_Env {
		table: make(map[string]SCMT),
		parent: parent,
	}
}
