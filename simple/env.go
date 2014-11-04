package simple


import (
	"github.com/JohnAnthony/goscm"
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
	env.BindSpecial("quote", scm_quote)
	env.BindSpecial("define", scm_define)
	env.BindSpecial("begin", scm_begin)
	env.BindSpecial("let", scm_let)
	env.BindSpecial("lambda", scm_lambda)
	env.BindSpecial("set!", scm_set_bang)
	return env
}

// WARNING! These procedures do no input validation, so feeding them incorrect
// input will have strange effects!

func scm_add(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	ret := 0
	for ; !args.IsNil(); args = args.Cdr.(*goscm.SCMT_Pair) {
		ret += args.Car.(*goscm.SCMT_Integer).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_multiply(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	ret := 1
	for ; !args.IsNil(); args = args.Cdr.(*goscm.SCMT_Pair) {
		ret *= args.Car.(*goscm.SCMT_Integer).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_subtract(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	ret := args.Car.(*goscm.SCMT_Integer).Value
	for args = args.Cdr.(*goscm.SCMT_Pair); !args.IsNil(); args = args.Cdr.(*goscm.SCMT_Pair) {
		ret -= args.Car.(*goscm.SCMT_Integer).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_divide(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	ret := args.Car.(*goscm.SCMT_Integer).Value
	for args = args.Cdr.(*goscm.SCMT_Pair); !args.IsNil(); args = args.Cdr.(*goscm.SCMT_Pair) {
		ret /= args.Car.(*goscm.SCMT_Integer).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_car(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	return args.Car.(*goscm.SCMT_Pair).Car, nil
}

func scm_cdr(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	return args.Car.(*goscm.SCMT_Pair).Cdr, nil
}

func scm_cons(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	return goscm.Cons(args.Car, args.Cdr.(*goscm.SCMT_Pair).Car), nil
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
	tobind, err := args.Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	env.Add(args.Car.(*goscm.SCMT_Symbol), tobind)
	return args.Car.(*goscm.SCMT_Symbol), nil
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
	return goscm.SCMT_Nil, env.Set(args.Car.(*goscm.SCMT_Symbol), args.Cdr.(*goscm.SCMT_Pair).Car)
}
