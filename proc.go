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
	var symb_lst SCMT

	argenv := EnvEmpty(p.env)
	arg_lst, err := MapEval(args, env)
	if err != nil {
		return SCMT_Nil, err
	}

	symb_lst = p.args
	for {
		if arg_lst == SCMT_Nil && symb_lst == SCMT_Nil {
			// A natural end to our zipping lists
			break
		} else if arg_lst == SCMT_Nil {
			// We ran out of symbols to attach to
			return SCMT_Nil, errors.New("Too few arguments")
		} else if symb_lst == SCMT_Nil {
			// We ran out of arguments to attach
			return SCMT_Nil, errors.New("Too many arguments")
		}

		// This handles the end case of a dotted argument list aka a variadic
		// function
		if reflect.TypeOf(symb_lst) == reflect.TypeOf(&SCMT_Symbol{}) {
			argenv.Add(symb_lst.(*SCMT_Symbol), arg_lst)
			break
		}
		
		argenv.Add(symb_lst.(*SCMT_Pair).Car.(*SCMT_Symbol), arg_lst.Car)
		arg_lst = arg_lst.Cdr.(*SCMT_Pair)
		symb_lst = symb_lst.(*SCMT_Pair).Cdr
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
