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

// Helpers

func (e *Environ) AddArgs(symbs *Pair, vals *Pair) error {
	// Both lists should terminate at the same time
	for symbs != SCM_Nil && vals != SCM_Nil {
		// Ran out of symbols first
		if symbs == SCM_Nil {
			return errors.New("Too many arguments")
		}

		// Ran out of arguments first
		if vals == SCM_Nil {
			return errors.New("Too few arguments")
		}
		
		symb, ok := symbs.Car.(*Symbol)
		if !ok {
			return errors.New("Argument list contained a non-symbol")
		}

		e.Add(symb, vals.Car)
		
		symbs, ok = symbs.Cdr.(*Pair)
		if !ok { // This is a dotted list
			symb, ok = symbs.Cdr.(*Symbol)
			if !ok {
				return errors.New("Argument list dotted with non-symbol")
			}

			e.Add(symb, vals.Cdr)
			break
		}
	}

	return nil
}
