package simple


import (
	"github.com/JohnAnthony/goscm"
	"reflect"
	"errors"
)

func Env() *goscm.SCMT_Env {
	env := goscm.EnvEmpty(nil)
	env.BindForeign("+", scm_add)
	env.BindForeign("-", scm_subtract)
	env.BindForeign("*", scm_multiply)
	env.BindForeign("/", scm_divide)
	env.BindForeign("car", scm_car)
	env.BindForeign("cdr", scm_cdr)
	env.BindForeign("cons", scm_cons)
	env.BindForeign("map", scm_map)
	env.BindForeign("apply", scm_apply)
	env.BindForeign("=", scm_numeq)
	env.BindSpecial("quote", scm_quote)
	env.BindSpecial("define", scm_define)
	env.BindSpecial("begin", scm_begin)
	env.BindSpecial("let", scm_let)
	env.BindSpecial("lambda", scm_lambda)
	env.BindSpecial("set!", scm_set_bang)
	env.BindSpecial("if", scm_if)
	return env
}

// WARNING! These procedures do no input validation, so feeding them incorrect
// input will have strange effects!

func scm_add(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	err := goscm.EnsureAll(args, reflect.TypeOf(&goscm.SCMT_Integer{}))
	if err != nil {
		return goscm.NilAndErr(err)
	}

	ret := 0
	for !args.IsNil() {
		ret += args.Car.(*goscm.SCMT_Integer).Value
		args, err = goscm.Cast_Pair(args.Cdr)
		if err != nil {
			return args, err
		}
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_multiply(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	err := goscm.EnsureAll(args, reflect.TypeOf(&goscm.SCMT_Integer{}))
	if err != nil {
		return goscm.NilAndErr(err)
	}

	ret := 1
	for !args.IsNil() {
		ret *= args.Car.(*goscm.SCMT_Integer).Value

		args, err = goscm.Cast_Pair(args.Cdr)
		if err != nil {
			return args, err
		}
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_subtract(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	err := goscm.EnsureAll(args, reflect.TypeOf(&goscm.SCMT_Integer{}))
	if err != nil {
		return goscm.NilAndErr(err)
	}

	ret := args.Car.(*goscm.SCMT_Integer).Value

	args, err = goscm.Cast_Pair(args.Cdr)
	if err != nil {
		return args, err
	}

	for !args.IsNil() {
		ret -= args.Car.(*goscm.SCMT_Integer).Value
		args, err = goscm.Cast_Pair(args.Cdr)
		if err != nil {
			return args, err
		}
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_divide(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	err := goscm.EnsureAll(args, reflect.TypeOf(&goscm.SCMT_Integer{}))
	if err != nil {
		return goscm.NilAndErr(err)
	}

	ret := args.Car.(*goscm.SCMT_Integer).Value

	args, err = goscm.Cast_Pair(args.Cdr)
	if err != nil {
		return args, err
	}

	for !args.IsNil() {
		ret /= args.Car.(*goscm.SCMT_Integer).Value
		args, err = goscm.Cast_Pair(args.Cdr)
		if err != nil {
			return args, err
		}
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_car(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	target, err := goscm.Cast_Pair(args.Car)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	return target.Car, nil
}

func scm_cdr(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	target, err := goscm.Cast_Pair(args.Car)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	return target.Cdr, nil
}

func scm_cons(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	first := args.Car

	second_cell, err := goscm.Cast_Pair(args.Cdr)
	if err != nil {
		return goscm.SCMT_Nil, errors.New("Not enough arguments")
	}

	second := second_cell.Car

	if second_cell.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	return goscm.Cons(first, second), nil
}

func scm_map(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	ret := goscm.SCMT_Nil
	f, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	target := args.Cdr.(*goscm.SCMT_Pair).Car.(*goscm.SCMT_Pair)
	for l := target; l != goscm.SCMT_Nil; l = l.Cdr.(*goscm.SCMT_Pair) {
		arg := goscm.Cons(l.Car, goscm.SCMT_Nil)
		applied, err := f.(goscm.SCMT_Func).Apply(arg, env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		ret = goscm.Cons(applied, ret)
	}
	return goscm.Reverse(ret), nil
}

func scm_apply(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	f := args.Car.(goscm.SCMT_Func)
	target := args.Cdr.(*goscm.SCMT_Pair).Car.(*goscm.SCMT_Pair)
	return f.Apply(target, env)
}

func scm_quote(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	return args.Car, nil
}

func scm_define(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	symb := args.Car
	switch reflect.TypeOf(symb) {
	case reflect.TypeOf(&goscm.SCMT_Symbol{}):
		tobind, err := args.Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		env.Add(symb.(*goscm.SCMT_Symbol), tobind)
		return symb, nil
	case reflect.TypeOf(&goscm.SCMT_Pair{}):
		name := symb.(*goscm.SCMT_Pair).Car.(*goscm.SCMT_Symbol)
		proc_args := symb.(*goscm.SCMT_Pair).Cdr
		proc_body := args.Cdr
		proc_tail := goscm.Cons(proc_args, proc_body)
		proc, err := scm_lambda(proc_tail, env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		env.Add(name, proc)
		return name, nil
	default:
		return goscm.SCMT_Nil,
		  errors.New("Attempting to define type: " + reflect.TypeOf(symb).String())
	}
}

func scm_begin(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	var result goscm.SCMT
	var err error

	for result = goscm.SCMT_Nil; args != goscm.SCMT_Nil; args = args.Cdr.(*goscm.SCMT_Pair) {
		result, err = args.Car.Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
	}
	return result, nil
}

func scm_let(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	body := args.Cdr.(*goscm.SCMT_Pair)
	newenv := goscm.EnvEmpty(env)
	
	for vars := args.Car.(*goscm.SCMT_Pair); vars != goscm.SCMT_Nil; vars = vars.Cdr.(*goscm.SCMT_Pair) {
		symb := vars.Car.(*goscm.SCMT_Pair).Car.(*goscm.SCMT_Symbol)
		val, err := vars.Car.(*goscm.SCMT_Pair).Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		newenv.Add(symb, val)
	}

	return scm_begin(body, newenv)
}

func scm_lambda(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	return goscm.Make_Proc(args.Car.(*goscm.SCMT_Pair), args.Cdr.(*goscm.SCMT_Pair), env), nil
}

func scm_set_bang(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	symb := args.Car.(*goscm.SCMT_Symbol)
	val, err := args.Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	return symb, env.Set(symb, val)
}

func scm_numeq(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	err := goscm.EnsureAll(args, reflect.TypeOf(&goscm.SCMT_Integer{}))
	if err != nil {
		return goscm.NilAndErr(err)
	}

	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}

	base := args.Car.(*goscm.SCMT_Integer).Value

	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		return goscm.WrongType()
	}
	args = args.Cdr.(*goscm.SCMT_Pair)
	
	for args != goscm.SCMT_Nil {
		if args.Car.(*goscm.SCMT_Integer).Value != base {
			return goscm.Make_SCMT(false), nil
		}

		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
			return goscm.WrongType()
		}
		args = args.Cdr.(*goscm.SCMT_Pair)
	}
	
	return goscm.Make_SCMT(true), nil
}

func scm_if(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	predicate, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		return goscm.WrongType()
	}
	
	if reflect.TypeOf(predicate) == reflect.TypeOf(&goscm.SCMT_Bool{}) &&
		predicate.(*goscm.SCMT_Bool).Value == false {

		if reflect.TypeOf(args.Cdr.(*goscm.SCMT_Pair).Cdr) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
			return goscm.WrongType()
		}
		return args.Cdr.(*goscm.SCMT_Pair).Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
	}
	
	return args.Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
}
