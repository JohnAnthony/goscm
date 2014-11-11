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
//TCO_TOP:
	for {
		newenv := EnvEmpty(env)
		newenv.AddArgs(p.args, args)
		body, err := p.body.ToSlice()
		if err != nil {
			return SCMT_Nil, err
		}

		for i := 0; i < len(body); i++ {
			body[i], err = body[i].Eval(newenv)
			if err != nil {
				return SCMT_Nil, err
			}

			if i == len(body) - 1 { // The last value, return its result
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
