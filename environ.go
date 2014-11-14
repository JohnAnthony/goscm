package goscm

import (
	"errors"
	"reflect"
)

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
	if ret != nil {
		return ret, nil
	}
	if env.parent == nil {
		return nil, errors.New("Unbound variable: " + symb.String())
	}
	return env.parent.Find(symb)
}

func (env *Environ) AddArgs(symbs *Pair, vals *Pair) error {
	for {
		if symbs == SCMT_Nil && vals == SCMT_Nil {
			// Natural end
			return nil
		}

		if symbs == SCMT_Nil {
			return errors.New("Too many arguments")
		}

		if vals == SCMT_Nil {
			return errors.New("Too few arguments")
		}
		
		if reflect.TypeOf(symbs.Car) != reflect.TypeOf(&Symbol{}) {
			return errors.New("Non-symbol passed as identifier")
		}

		env.Add(symbs.Car.(*Symbol), vals.Car)
		
		if reflect.TypeOf(symbs.Cdr) != reflect.TypeOf(&Pair{}) {
			env.Add(symbs.Cdr.(*Symbol), vals.Cdr)
			return nil
		}
		
		if reflect.TypeOf(vals.Cdr) != reflect.TypeOf(&Pair{}) {
			return errors.New("Dotted argument list")
		}
		
		symbs = symbs.Cdr.(*Pair)
		vals = vals.Cdr.(*Pair)
	}
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

func (env *Environ) BindForeign(name string, f func (*Pair, *Environ) (SCMT, error)) {
	env.Add(Make_Symbol(name), Make_Foreign(f))
}

func (env *Environ) BindSpecial(name string, f func (*Pair, *Environ) (SCMT, error)) {
	env.Add(Make_Symbol(name), Make_Special(f))
}

func EnvEmpty(parent *Environ) *Environ {
	return &Environ {
		table: make(map[string]SCMT),
		parent: parent,
	}
}
