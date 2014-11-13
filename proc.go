package goscm

import (
	"errors"
	"reflect"
)

type SCMT_Proc struct {
	args *SCMT_Pair
	body *SCMT_Pair
	env *SCMT_Env
}

func (p *SCMT_Proc) Eval(*SCMT_Env) (SCMT, error) {
	return p, nil
}

func (*SCMT_Proc) String() string {
	return "#<procedure>"
}

func (p *SCMT_Proc) Apply(args *SCMT_Pair, env *SCMT_Env) (SCMT, error) {
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
				for reflect.TypeOf(body[i]) == reflect.TypeOf(&SCMT_Pair{}) {
					tail_pair := body[i].(*SCMT_Pair)
					
					if reflect.TypeOf(tail_pair.Car) != reflect.TypeOf(&SCMT_Symbol{}) {
						break
					}

					tail_func, err := tail_pair.Car.(*SCMT_Symbol).Eval(newenv)
					if err != nil { return SCMT_Nil, err }
					
					if reflect.TypeOf(tail_func) == reflect.TypeOf(&SCMT_Special{}) {
						if tail_func.(*SCMT_Special).Expand != nil {
							// TODO: Check that body[i].Cdr is a *SCMT_Pair
							body[i], err = tail_func.(*SCMT_Special).Expand(body[i].(*SCMT_Pair).Cdr.(*SCMT_Pair), newenv)
							if err != nil { return SCMT_Nil, err }
							continue
						}
					}

					if tail_func == p {
						env = newenv
						args = body[i].(*SCMT_Pair).Cdr.(*SCMT_Pair)
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

func Make_Proc(args *SCMT_Pair, body *SCMT_Pair, env *SCMT_Env) (*SCMT_Proc, error) {
	for a := args; a != SCMT_Nil; a = a.Cdr.(*SCMT_Pair) {
		if reflect.TypeOf(a.Cdr) == reflect.TypeOf(&SCMT_Symbol{}) {
			break
		}
		if reflect.TypeOf(a.Cdr) != reflect.TypeOf(&SCMT_Pair{}) {
			return nil, errors.New("Non-symbol used to create procedure")
		}
	}

	return &SCMT_Proc {
		args: args,
		body: body,
		env: env,
	}, nil
}
