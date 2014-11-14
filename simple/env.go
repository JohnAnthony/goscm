package simple


import (
	"github.com/JohnAnthony/goscm"
	"reflect"
	"errors"
)

func Env() *goscm.Environ {
	env := goscm.EnvEmpty(nil)
	env.BindForeign("*", scm_multiply)
	env.BindForeign("+", scm_add)
	env.BindForeign("-", scm_subtract)
	env.BindForeign("/", scm_divide)
	env.BindForeign("<", scm_lt)
	env.BindForeign("<=", scm_le)
	env.BindForeign("=", scm_numeq)
	env.BindForeign(">", scm_gt)
	env.BindForeign(">=", scm_ge)
	env.BindForeign("apply", scm_apply)
	env.BindForeign("car", scm_car)
	env.BindForeign("cdr", scm_cdr)
	env.BindForeign("cons", scm_cons)
	env.BindForeign("list", scm_list)
	env.BindForeign("map", scm_map)
	env.BindForeign("not", scm_not)
	env.BindForeign("reverse", scm_reverse)
	env.BindSpecial("and", scm_and)
	env.BindSpecial("begin", scm_begin)
	env.BindSpecial("define", scm_define)
	env.BindSpecial("lambda", scm_lambda)
	env.BindSpecial("let", scm_let)
	env.BindSpecial("or", scm_or)
	env.BindSpecial("quote", scm_quote)
	env.BindSpecial("set!", scm_set_bang)
	env.Add(goscm.Make_Symbol("if"), goscm.Make_SpecialExpandable(scm_if, scm_if_expand))
	return env
}

func scm_add(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	ret := 0
	for i := 0; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer type")
		}
		ret += argss[i].(*goscm.PlainInt).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_multiply(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	ret := 1
	for i := 0; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer type")
		}
		ret *= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_subtract(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if len(argss) == 0 {
		return goscm.SCMT_Nil, errors.New("Expected at least one argument")
	}
	
	if reflect.TypeOf(argss[0]) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer type")
	}
	ret := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer type")
		}
		ret -= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_divide(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if len(argss) == 0 {
		return goscm.SCMT_Nil, errors.New("Expected at least one argument")
	}
	
	if reflect.TypeOf(argss[0]) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer type")
	}
	ret := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer type")
		}
		ret /= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.Make_SCMT(ret), nil
}

func scm_car(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

func scm_cdr(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

func scm_cons(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

func scm_map(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Expected pair")
	}

	tail, err := argss[1].(*goscm.Pair).ToSlice()
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

func scm_apply(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Non-list to apply to")
	}
	target := argss[1].(*goscm.Pair)
	return f.Apply(target, env)
}

func scm_quote(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too few arguments")
	}
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}
	return args.Car, nil
}

func scm_define(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	symb := args.Car
	switch reflect.TypeOf(symb) {
	case reflect.TypeOf(&goscm.Symbol{}):
		tobind, err := args.Cdr.(*goscm.Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		env.Add(symb.(*goscm.Symbol), tobind)
		return symb, nil
	case reflect.TypeOf(&goscm.Pair{}):
		name := symb.(*goscm.Pair).Car.(*goscm.Symbol)
		proc_args := symb.(*goscm.Pair).Cdr
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

func scm_begin(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
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

func scm_let(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	body := args.Cdr.(*goscm.Pair)
	newenv := goscm.EnvEmpty(env)
	
	for vars := args.Car.(*goscm.Pair); vars != goscm.SCMT_Nil; vars = vars.Cdr.(*goscm.Pair) {
		symb := vars.Car.(*goscm.Pair).Car.(*goscm.Symbol)
		val, err := vars.Car.(*goscm.Pair).Cdr.(*goscm.Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCMT_Nil, err
		}
		newenv.Add(symb, val)
	}

	return scm_begin(body, newenv)
}

func scm_lambda(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	proc, err := goscm.Make_Proc(
		args.Car.(*goscm.Pair),
		args.Cdr.(*goscm.Pair),
		env,
	)

	if err != nil {
		return goscm.SCMT_Nil, err
	}

	return proc, nil
}

func scm_set_bang(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	symb := args.Car.(*goscm.Symbol)
	val, err := args.Cdr.(*goscm.Pair).Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	return symb, env.Set(symb, val)
}

func scm_numeq(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}

	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	base := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if argss[i].(*goscm.PlainInt).Value != base {
			return goscm.Make_SCMT(false), nil
		}
	}

	return goscm.Make_SCMT(true), nil
}

func scm_lt(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if args.Cdr == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCMT_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if args.Car.(*goscm.PlainInt).Value >= next.(*goscm.PlainInt).Value {
			return goscm.Make_SCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.Make_SCMT(true), nil
}

func scm_gt(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if args.Cdr == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCMT_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if args.Car.(*goscm.PlainInt).Value <= next.(*goscm.PlainInt).Value {
			return goscm.Make_SCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.Make_SCMT(true), nil
}

func scm_le(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if args.Cdr == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCMT_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if !(args.Car.(*goscm.PlainInt).Value <= next.(*goscm.PlainInt).Value) {
			return goscm.Make_SCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.Make_SCMT(true), nil
}

func scm_ge(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if args.Cdr == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCMT_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCMT_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCMT_Nil, errors.New("Expected integer argument")
		}

		if !(args.Car.(*goscm.PlainInt).Value >= next.(*goscm.PlainInt).Value) {
			return goscm.Make_SCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.Make_SCMT(true), nil
}

func scm_and(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(true), nil
	}

	result, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	if !goscm.IsTrue(result) {
		return goscm.Make_SCMT(false), nil
	}

	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
	}

	return scm_and(args.Cdr.(*goscm.Pair), env)
}

func scm_or(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.Make_SCMT(false), nil
	}

	result, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	if goscm.IsTrue(result) {
		return goscm.Make_SCMT(true), nil
	}

	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Expected nil-terminated list")
	}

	return scm_or(args.Cdr.(*goscm.Pair), env)
}

func scm_if(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	expanded, err := scm_if_expand(args, env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	return expanded.Eval(env)
}

func scm_if_expand(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	predicate, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCMT_Nil, err
	}
	
	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Wrong type")
	}
	
	if reflect.TypeOf(predicate) == reflect.TypeOf(&goscm.Boolean{}) &&
		predicate.(*goscm.Boolean).Value == false {

		if reflect.TypeOf(args.Cdr.(*goscm.Pair).Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCMT_Nil, errors.New("Wrong type")
		}
		return args.Cdr.(*goscm.Pair).Cdr.(*goscm.Pair).Car, nil
	}
	
	return args.Cdr.(*goscm.Pair).Car, nil
}

func scm_not(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Too many arguments")
	}
	
	if reflect.TypeOf(args.Car) == reflect.TypeOf(&goscm.Boolean{}) {
		return goscm.Make_SCMT(!args.Car.(*goscm.Boolean).Value), nil
	}

	return goscm.Make_SCMT(false), nil
}

func scm_reverse(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Expected exactly one argument")
	}
	if args.Cdr != goscm.SCMT_Nil {
		return goscm.SCMT_Nil, errors.New("Expected exactly one argument")
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCMT_Nil, errors.New("Expected list type")
	}

	list, err := args.Car.(*goscm.Pair).ToSlice()
	if err != nil {
		return goscm.SCMT_Nil, err
	}

	ret := goscm.SCMT_Nil
	for i := 0; i < len(list); i++ {
		ret = goscm.Cons(list[i], ret)
	}
	return ret, nil
}

func scm_list(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	return args, nil
}
