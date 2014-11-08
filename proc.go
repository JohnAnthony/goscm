package goscm

import "errors"

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
	argenv := EnvEmpty(p.env)

	arg := args
	symb := p.args
	for {
		if arg == SCMT_Nil && symb == SCMT_Nil {
			break
		}
		if arg == SCMT_Nil {
			// We ran out of symbols to attach to
			return SCMT_Nil, errors.New("Too few arguments")
		}
		if symb == SCMT_Nil {
			return SCMT_Nil, errors.New("Too many arguments")
		}

		val, err := arg.Car.Eval(env)
		if err != nil {
			return SCMT_Nil, err
		}

		argenv.Add(symb.Car.(*SCMT_Symbol), val)
		arg = arg.Cdr.(*SCMT_Pair)
		symb = symb.Cdr.(*SCMT_Pair)
		
	}

	var result SCMT
	for expr := p.body; expr != SCMT_Nil; expr = expr.Cdr.(*SCMT_Pair) {
		result, err = expr.Car.Eval(argenv)
		if err != nil {
			return SCMT_Nil, err
		}
	}

	return result, nil
}

func Make_Proc(args *SCMT_Pair, body *SCMT_Pair, env *SCMT_Env) *SCMT_Proc {
	return &SCMT_Proc {
		args: args,
		body: body,
		env: env,
	}
}
