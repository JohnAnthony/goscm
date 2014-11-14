package goscm

import (
	"errors"
	"reflect"
)

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
		args, err = MapEval(args, env)
		if err != nil {	return SCMT_Nil, err }

		newenv := EnvEmpty(env)
		newenv.AddArgs(p.args, args)
		body, err := p.body.ToSlice()
		if err != nil {	return SCMT_Nil, err }

		for i := 0; i < len(body); i++ {
			
			// This is the last call in the body. We might need to TCO
			if i == len(body) - 1 {
				for reflect.TypeOf(body[i]) == reflect.TypeOf(&Pair{}) {
					tail_pair := body[i].(*Pair)
					
					if reflect.TypeOf(tail_pair.Car) != reflect.TypeOf(&Symbol{}) {
						break
					}

					tail_func, err := tail_pair.Car.(*Symbol).Eval(newenv)
					if err != nil { return SCMT_Nil, err }
					
					if reflect.TypeOf(tail_func) == reflect.TypeOf(&SCMT_Special{}) {
						if tail_func.(*SCMT_Special).Expand != nil {
							// TODO: Check that body[i].Cdr is a *Pair
							body[i], err = tail_func.(*SCMT_Special).Expand(body[i].(*Pair).Cdr.(*Pair), newenv)
							if err != nil { return SCMT_Nil, err }
							continue
						}
					}

					if tail_func == p {
						env = newenv
						args = body[i].(*Pair).Cdr.(*Pair)
						continue TCO_TOP
					}

					break
				}
			}

			// This is not the last expression in the body, continue normally
			body[i], err = body[i].Eval(newenv)
			if err != nil { return SCMT_Nil, err }
			
			if i == len(body) - 1 {
				return body[i], nil
			}
		}

		break // Should never get here, but if we do we want to get out and error
	}
	
	return SCMT_Nil, errors.New("Execution flow got somewhere it shouldn't")
}

func Make_Proc(args *Pair, body *Pair, env *Environ) (*Proc, error) {
	for a := args; a != SCMT_Nil; a = a.Cdr.(*Pair) {
		if reflect.TypeOf(a.Cdr) == reflect.TypeOf(&Symbol{}) {
			break
		}
		if reflect.TypeOf(a.Cdr) != reflect.TypeOf(&Pair{}) {
			return nil, errors.New("Non-symbol used to create procedure")
		}
	}

	return &Proc {
		args: args,
		body: body,
		env: env,
	}, nil
}
