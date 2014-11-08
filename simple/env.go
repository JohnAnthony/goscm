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
	env.BindForeign("not", scm_not)
	env.BindSpecial("quote", scm_quote)
	env.BindSpecial("define", scm_define)
	env.BindSpecial("begin", scm_begin)
	env.BindSpecial("let", scm_let)
	env.BindSpecial("lambda", scm_lambda)
	env.BindSpecial("set!", scm_set_bang)
	env.BindSpecial("if", scm_if)
	return env
}

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
	if args == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}

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
	if args == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}

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
		return goscm.SCMT_Nil, err
	}
	if second_cell == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}

	second := second_cell.Car

	if second_cell.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	return goscm.Cons(first, second), nil
}

func scm_map(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	if len(argss) < 2 {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}
	if len(argss) > 2 {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	proc, ok := argss[0].(goscm.SCMT_Func)
	if !ok {
		return goscm.SCMT_Nil, errors.New("Non-function in first position")
	}

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		return goscm.SCMT_Nil, errors.New("Expected pair")
	}

	tail, err := argss[1].(*goscm.SCMT_Pair).ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	ret := goscm.SCMT_Nil
	for i := len(tail) - 1; i >= 0; i-- {
		result, err := proc.Apply(goscm.Make_List(tail[i]), env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		ret = goscm.Cons(result, ret)
	}
	
	return ret, nil
}

func scm_apply(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	if len(argss) < 2 {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}
	if len(argss) > 3 {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}

	f, ok := argss[0].(goscm.SCMT_Func)
	if !ok {
		return goscm.SCMT_Nil, errors.New("Attempt to apply non-function type")
	}

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.SCMT_Pair{}) {
		return goscm.SCMT_Nil, errors.New("Non-list to apply to")
	}
	target := argss[1].(*goscm.SCMT_Pair)
	return f.Apply(target, env)
}

func scm_quote(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}
	return args.Car, nil
}

func scm_define(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	// Unsafe
	// No argument number checking
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

	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	for i := 0; i < len(argss); i++ {
		result, err = argss[i].Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
	}

	return result, nil
}

func scm_let(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	// Unsafe
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
	// Unsafe
	// No argument number checking
	proc, err := goscm.Make_Proc(
		args.Car.(*goscm.SCMT_Pair),
		args.Cdr.(*goscm.SCMT_Pair),
		env,
	)

	if err != nil {
		return goscm.SCMT_Nil, err
	}

	return proc, nil
}

func scm_set_bang(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	// Unsafe
	// No argument number checking
	symb := args.Car.(*goscm.SCMT_Symbol)
	val, err := args.Cdr.(*goscm.SCMT_Pair).Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	return symb, env.Set(symb, val)
}

func scm_numeq(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}

	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	base := argss[0].(*goscm.SCMT_Integer).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.SCMT_Integer{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if argss[i].(*goscm.SCMT_Integer).Value != base {
			return goscm.Make_SCMT(false), nil
		}
	}

	return goscm.Make_SCMT(true), nil
}

func scm_if(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	// Unsafe
	// No argument number checking
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

func scm_not(args *goscm.SCMT_Pair, env *goscm.SCMT_Env) (goscm.SCMT, error) {
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}
	
	if reflect.TypeOf(args.Car) == reflect.TypeOf(&goscm.SCMT_Bool{}) {
		return goscm.Make_SCMT(!args.Car.(*goscm.SCMT_Bool).Value), nil
	}

	return goscm.Make_SCMT(false), nil
}

// TODO: cond
// TODO: <=
// TODO: >=
// TODO: < 
// TODO: > 
// TODO: and
// TODO: or
