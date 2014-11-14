package goscm

import "errors"

type Environ struct {
	table map[string]SCMT
	parent *Environ
}

func (env *Environ) Eval(*Environ) (SCMT, error) {
	return env, nil
}

func (*Environ) String() string {
	return "#<environment>"
}

func (env *Environ) Add(symb *Symbol, val SCMT) {
	env.table[symb.Value] = val
}

func (env *Environ) Find(symb *Symbol) (SCMT, error) {
	ret := env.table[symb.Value]
	if ret != nil { // We found something
		return ret, nil
	}
	if env.parent == nil {
		return nil, errors.New("Unbound variable: " + symb.String())
	}
	return env.parent.Find(symb)
}

func (env *Environ) Set(symb *Symbol, val SCMT) error {
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

func EnvEmpty(parent *Environ) *Environ {
	return &Environ {
		table: make(map[string]SCMT),
		parent: parent,
	}
}
