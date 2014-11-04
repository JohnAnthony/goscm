package goscm

type SCMT_Proc struct {
	args *SCMT_Pair
	body *SCMT_Pair
	env *SCMT_Env
}

func (p *SCMT_Proc) Eval(*SCMT_Env) SCMT {
	return p
}

func (*SCMT_Proc) String() string {
	return "#<procedure>"
}

func (p *SCMT_Proc) Apply(args *SCMT_Pair, env *SCMT_Env) SCMT {
	argenv := EnvEmpty(p.env)

	arg := args
	symb := p.args
	for arg != SCMT_Nil && symb != SCMT_Nil {
		argenv.Add(symb.Car.(*SCMT_Symbol), arg.Car.Eval(env))
		arg = arg.Cdr.(*SCMT_Pair)
		symb = symb.Cdr.(*SCMT_Pair)
		
		if arg == SCMT_Nil && symb == SCMT_Nil {
			break
		}
	}

	var result SCMT
	for expr := p.body; expr != SCMT_Nil; expr = expr.Cdr.(*SCMT_Pair) {
		result = expr.Car.Eval(argenv)
	}

	return result
}

func Make_Proc(args *SCMT_Pair, body *SCMT_Pair, env *SCMT_Env) *SCMT_Proc {
	return &SCMT_Proc {
		args: args,
		body: body,
		env: env,
	}
}
