package goscm

type SCMT_Special struct {
	function func (*Pair, *Environ) (SCMT, error)
	Expand func (*Pair, *Environ) (SCMT, error) // Used for TCO
}

func (sp *SCMT_Special) Eval(*Environ) (SCMT, error) {
	return sp, nil
}

func (*SCMT_Special) String() string {
	return "#<special form>"
}

func (sp *SCMT_Special) Apply(args *Pair, env *Environ) (SCMT, error) {
	return sp.function(args, env)
}

func NewSpecial(f func (*Pair, *Environ) (SCMT, error)) *SCMT_Special {
	return &SCMT_Special {
		function: f,
		Expand: nil,
	}
}

func NewSpecialExpandable(f func (*Pair, *Environ) (SCMT, error), ex func (*Pair, *Environ) (SCMT, error)) *SCMT_Special {
	return &SCMT_Special {
		function: f,
		Expand: ex,
	}
}
