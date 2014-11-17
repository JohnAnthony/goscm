package goscm

import "errors"

type Proc struct {
	args *Pair
	body *Pair
	env *Environ
}

func (p *Proc) Eval(*Environ) (SCMT, error) {
	return p, nil
}

func (*Proc) String() string {
	return "#<procedure>"
}

func (p *Proc) Apply(args *Pair, env *Environ) (SCMT, error) {
	var err error

TCO_TOP:
	for {
		args, err = args.MapEval(env)
		if err != nil {	return SCM_Nil, err }

		newenv := NewEnv(env)
		newenv.AddArgs(p.args, args)
		body, err := p.body.ToSlice()
		if err != nil {	return SCM_Nil, err }

		for i := 0; i < len(body); i++ {
			// This is the last call in the body. We might need to TCO
			if i == len(body) - 1 {
				// We use a for loop here because later we may be replacing
				// body[i] as we partiall evaluate special forms that allow for
				// TCO. If we get another special form that may allow TCO we
				// have to jump back to this point and do it all again. Repeat
				// until we're out of expandable special forms or TCO is
				// confirmed.
				for {
					// If it's not a pair, it can't be a TCO candidate
					tail_pair, ok := body[i].(*Pair)
					if !ok { break }
					
					// If it's not a symbol, it can't be a TCO candidate
					symb, ok := tail_pair.Car.(*Symbol)
					if !ok { break }

					// Eval the symbol to get what's actually in TCO position
					tail_func, err := symb.Eval(newenv)
					if err != nil { return SCM_Nil, err }
					
					// Some special forms can provide TCO, but we have to expand
					// them first to find out.
					if spesh, ok := tail_func.(*Special); ok {
						if spesh.Expand != nil {
							spesh_args, ok := tail_pair.Cdr.(*Pair)
							if !ok { break }

							body[i], err = spesh.Expand(spesh_args, newenv)
							if err != nil { return SCM_Nil, err }
							
							// We've replaced body[i] with something else. It
							// might be suitable for TCO, but we need to check
							// everything we just checked all over again. Back
							// to the top.
							continue
						}
					}

					// If the TCO position Proc matches we have TCO. Stat from
					// the top of the function again with arguments swapped to
					// the new ones.
					if tail_func == p {
						env = newenv
						args = body[i].(*Pair).Cdr.(*Pair)
						continue TCO_TOP
					}

					// If we get here, TCo has failed. We have some function
					// other than p in tail position.
					break
				}
			}

			// This is not the last expression in the body, continue normally
			body[i], err = body[i].Eval(newenv)
			if err != nil { return SCM_Nil, err }
			
			if i == len(body) - 1 {
				return body[i], nil
			}
		}

		break // Should never get here, but if we do we want to error out
	}
	
	return SCM_Nil, errors.New("Execution flow got somewhere it shouldn't")
}

func NewProc(args *Pair, body *Pair, env *Environ) (*Proc, error) {
	var ok bool

	// Check that all of our arguments are symbols. Also accept a cdr that is a
	// symbol because that is a dotted list used for variadic functions.
	a := args
	for a != SCM_Nil {
		_, ok = a.Cdr.(*Symbol)
		if ok { break }

		a, ok = a.Cdr.(*Pair)
		if !ok {
			return nil, errors.New("Non-symbol used to create procedure")
		}
	}

	return &Proc {
		args: args,
		body: body,
		env: env,
	}, nil
}
