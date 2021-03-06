package goscm

import "errors"

// Special comes in two flavours - expanding and non-expanding. Expanding
// Specials are used to provide TCO, but must be expanded before they can
// be Applied in other contexts. Other Special forms can jsut be applied
// straight away.

// If you want to make a TCOExpand-type special form, the function you give
// should actually return something else to Eval.
type Special struct {
	Function func (*Pair, *Environ) (SCMT, error)
	TCOExpand bool
}

func (sp *Special) Eval(*Environ) (SCMT, error) {
	return sp, nil
}

func (*Special) String() string {
	return "#<special form>"
}

func (sp *Special) Apply(args *Pair, env *Environ) (SCMT, error) {
	// Expanding functions
	if sp.TCOExpand {
		expanded, err := sp.Expand(args, env)
		if err != nil { return SCM_Nil, err }
		
		return expanded.Eval(env)
	}

	// Non-expanding functions just get called
	return sp.Function(args, env)
}

func (sp *Special) Expand(args *Pair, env *Environ) (SCMT, error) {
	if sp.TCOExpand {
		return sp.Function(args, env)
	}
	return SCM_Nil, errors.New("Attempt to expand non-expanding special form")
}

func NewSpecial(f func (*Pair, *Environ) (SCMT, error)) *Special {
	return &Special {
		Function: f,
		TCOExpand: false,
	}
}

func NewSpecialTCO(f func (*Pair, *Environ) (SCMT, error)) *Special {
	return &Special {
		Function: f,
		TCOExpand: true,
	}
}
