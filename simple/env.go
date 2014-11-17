package simple

import (
	"github.com/JohnAnthony/goscm"
	"reflect"
	"errors"
)

func Env() *goscm.Environ {
	env := goscm.EnvEmpty(nil)
	env.Add(goscm.NewSymbol("*"), goscm.NewForeign(scm_multiply))
	env.Add(goscm.NewSymbol("+"), goscm.NewForeign(scm_add))
	env.Add(goscm.NewSymbol("-"), goscm.NewForeign(scm_subtract))
	env.Add(goscm.NewSymbol("/"), goscm.NewForeign(scm_divide))
	env.Add(goscm.NewSymbol("<"), goscm.NewForeign(scm_lt))
	env.Add(goscm.NewSymbol("<="), goscm.NewForeign(scm_le))
	env.Add(goscm.NewSymbol("="), goscm.NewForeign(scm_numeq))
	env.Add(goscm.NewSymbol(">"), goscm.NewForeign(scm_gt))
	env.Add(goscm.NewSymbol(">="), goscm.NewForeign(scm_ge))
	env.Add(goscm.NewSymbol("apply"), goscm.NewForeign(scm_apply))
	env.Add(goscm.NewSymbol("car"), goscm.NewForeign(scm_car))
	env.Add(goscm.NewSymbol("cdr"), goscm.NewForeign(scm_cdr))
	env.Add(goscm.NewSymbol("cons"), goscm.NewForeign(scm_cons))
	env.Add(goscm.NewSymbol("list"), goscm.NewForeign(scm_list))
	env.Add(goscm.NewSymbol("map"), goscm.NewForeign(scm_map))
	env.Add(goscm.NewSymbol("not"), goscm.NewForeign(scm_not))
	env.Add(goscm.NewSymbol("reverse"), goscm.NewForeign(scm_reverse))
	env.Add(goscm.NewSymbol("and"), goscm.NewSpecial(scm_and))
	env.Add(goscm.NewSymbol("begin"), goscm.NewSpecial(scm_begin))
	env.Add(goscm.NewSymbol("define"), goscm.NewSpecial(scm_define))
	env.Add(goscm.NewSymbol("lambda"), goscm.NewSpecial(scm_lambda))
	env.Add(goscm.NewSymbol("let"), goscm.NewSpecial(scm_let))
	env.Add(goscm.NewSymbol("or"), goscm.NewSpecial(scm_or))
	env.Add(goscm.NewSymbol("quote"), goscm.NewSpecial(scm_quote))
	env.Add(goscm.NewSymbol("set!"), goscm.NewSpecial(scm_set_bang))
	env.Add(goscm.NewSymbol("if"), goscm.NewSpecialTCO(scm_if))
	env.Add(goscm.NewSymbol("cond"), goscm.NewSpecialTCO(scm_cond))
	return env
}

func scm_add(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	ret := 0
	for i := 0; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer type")
		}
		ret += argss[i].(*goscm.PlainInt).Value
	}
	return goscm.NewSCMT(ret), nil
}

func scm_multiply(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	ret := 1
	for i := 0; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer type")
		}
		ret *= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.NewSCMT(ret), nil
}

func scm_subtract(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}
	
	if len(argss) == 0 {
		return goscm.SCM_Nil, errors.New("Expected at least one argument")
	}
	
	if reflect.TypeOf(argss[0]) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer type")
	}
	ret := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer type")
		}
		ret -= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.NewSCMT(ret), nil
}

func scm_divide(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}
	
	if len(argss) == 0 {
		return goscm.SCM_Nil, errors.New("Expected at least one argument")
	}
	
	if reflect.TypeOf(argss[0]) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer type")
	}
	ret := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer type")
		}
		ret /= argss[i].(*goscm.PlainInt).Value
	}
	return goscm.NewSCMT(ret), nil
}

func scm_car(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}

	target, ok := args.Car.(*goscm.Pair)
	if !ok {
		return goscm.SCM_Nil, errors.New("Attempt to take Car of non-list")
	}
	
	if args.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}

	return target.Car, nil
}

func scm_cdr(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}

	target, ok := args.Car.(*goscm.Pair)
	if !ok {
		return goscm.SCM_Nil, errors.New("Attempt to take cdr of non-list")
	}
	
	if args.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}

	return target.Cdr, nil
}

func scm_cons(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	first := args.Car

	second_cell, ok := args.Cdr.(*goscm.Pair)
	if !ok {
		return goscm.SCM_Nil, errors.New("Weird (unexpected) dotted list")
	}

	if second_cell == goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}

	second := second_cell.Car

	if second_cell.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}

	return goscm.Cons(first, second), nil
}

func scm_map(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	if len(argss) < 2 {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}
	if len(argss) > 2 {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}

	proc, ok := argss[0].(goscm.Func)
	if !ok {
		return goscm.SCM_Nil, errors.New("Non-function in first position")
	}

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Expected pair")
	}

	tail, err := argss[1].(*goscm.Pair).ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	ret := goscm.SCM_Nil
	for i := len(tail) - 1; i >= 0; i-- {
		result, err := proc.Apply(goscm.NewList(tail[i]), env)
		if err != nil {
			return goscm.SCM_Nil, err
		}
		ret = goscm.Cons(result, ret)
	}
	
	return ret, nil
}

func scm_apply(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}
	if len(argss) < 2 {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}
	if len(argss) > 3 {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}

	f, ok := argss[0].(goscm.Func)
	if !ok {
		return goscm.SCM_Nil, errors.New("Attempt to apply non-function type")
	}

	if reflect.TypeOf(argss[1]) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Non-list to apply to")
	}
	target := argss[1].(*goscm.Pair)
	return f.Apply(target, env)
}

func scm_quote(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too few arguments")
	}
	if args.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}
	return args.Car, nil
}

func scm_define(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	symb := args.Car
	switch reflect.TypeOf(symb) {
	case reflect.TypeOf(&goscm.Symbol{}):
		tobind, err := args.Cdr.(*goscm.Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCM_Nil, err
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
			return goscm.SCM_Nil, err
		}
		env.Add(name, proc)
		return name, nil
	default:
		return goscm.SCM_Nil,
		  errors.New("Attempting to define type: " + reflect.TypeOf(symb).String())
	}
}

func scm_begin(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	var result goscm.SCMT

	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	for i := 0; i < len(argss); i++ {
		result, err = argss[i].Eval(env)
		if err != nil {
			return goscm.SCM_Nil, err
		}
	}

	return result, nil
}

func scm_let(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	body := args.Cdr.(*goscm.Pair)
	newenv := goscm.EnvEmpty(env)
	
	for vars := args.Car.(*goscm.Pair); vars != goscm.SCM_Nil; vars = vars.Cdr.(*goscm.Pair) {
		symb := vars.Car.(*goscm.Pair).Car.(*goscm.Symbol)
		val, err := vars.Car.(*goscm.Pair).Cdr.(*goscm.Pair).Car.Eval(env)
		if err != nil {
			return goscm.SCM_Nil, err
		}
		newenv.Add(symb, val)
	}

	return scm_begin(body, newenv)
}

func scm_lambda(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	proc, err := goscm.NewProc(
		args.Car.(*goscm.Pair),
		args.Cdr.(*goscm.Pair),
		env,
	)

	if err != nil {
		return goscm.SCM_Nil, err
	}

	return proc, nil
}

func scm_set_bang(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	symb := args.Car.(*goscm.Symbol)
	val, err := args.Cdr.(*goscm.Pair).Car.Eval(env)
	if err != nil {
		return goscm.SCM_Nil, err
	}
	return symb, env.Set(symb, val)
}

func scm_numeq(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}

	argss, err := args.ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	base := argss[0].(*goscm.PlainInt).Value

	for i := 1; i < len(argss); i++ {
		if reflect.TypeOf(argss[i]) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer argument")
		}

		if argss[i].(*goscm.PlainInt).Value != base {
			return goscm.NewSCMT(false), nil
		}
	}

	return goscm.NewSCMT(true), nil
}

func scm_lt(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if args.Cdr == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCM_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer argument")
		}

		if args.Car.(*goscm.PlainInt).Value >= next.(*goscm.PlainInt).Value {
			return goscm.NewSCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.NewSCMT(true), nil
}

func scm_gt(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if args.Cdr == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCM_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer argument")
		}

		if args.Car.(*goscm.PlainInt).Value <= next.(*goscm.PlainInt).Value {
			return goscm.NewSCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.NewSCMT(true), nil
}

func scm_le(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if args.Cdr == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCM_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer argument")
		}

		if !(args.Car.(*goscm.PlainInt).Value <= next.(*goscm.PlainInt).Value) {
			return goscm.NewSCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.NewSCMT(true), nil
}

func scm_ge(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if args.Cdr == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.PlainInt{}) {
		return goscm.SCM_Nil, errors.New("Expected integer argument")
	}

	for args.Cdr != goscm.SCM_Nil {
		if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
		}
		
		next := args.Cdr.(*goscm.Pair).Car
		if reflect.TypeOf(next) != reflect.TypeOf(&goscm.PlainInt{}) {
			return goscm.SCM_Nil, errors.New("Expected integer argument")
		}

		if !(args.Car.(*goscm.PlainInt).Value >= next.(*goscm.PlainInt).Value) {
			return goscm.NewSCMT(false), nil
		}

		args = args.Cdr.(*goscm.Pair)
	}

	return goscm.NewSCMT(true), nil
}

func scm_and(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(true), nil
	}

	result, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCM_Nil, err
	}

	if !goscm.SCMTrue(result) {
		return goscm.NewSCMT(false), nil
	}

	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
	}

	return scm_and(args.Cdr.(*goscm.Pair), env)
}

func scm_or(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.NewSCMT(false), nil
	}

	result, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCM_Nil, err
	}

	if goscm.SCMTrue(result) {
		return goscm.NewSCMT(true), nil
	}

	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Expected nil-terminated list")
	}

	return scm_or(args.Cdr.(*goscm.Pair), env)
}

func scm_if(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	predicate, err := args.Car.Eval(env)
	if err != nil {
		return goscm.SCM_Nil, err
	}
	
	if reflect.TypeOf(args.Cdr) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Wrong type")
	}
	
	if reflect.TypeOf(predicate) == reflect.TypeOf(&goscm.Boolean{}) &&
		predicate.(*goscm.Boolean).Value == false {

		if reflect.TypeOf(args.Cdr.(*goscm.Pair).Cdr) != reflect.TypeOf(&goscm.Pair{}) {
			return goscm.SCM_Nil, errors.New("Wrong type")
		}
		return args.Cdr.(*goscm.Pair).Cdr.(*goscm.Pair).Car, nil
	}
	
	return args.Cdr.(*goscm.Pair).Car, nil
}

func scm_cond(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	for args != goscm.SCM_Nil {
		condition, ok := args.Car.(*goscm.Pair)
		if !ok {
			return goscm.SCM_Nil, errors.New("Expected list")
		}
		
		pred, err := condition.Car.Eval(env)
		if err != nil { return goscm.SCM_Nil, err }

		if goscm.SCMTrue(pred) {
			consequent_cell, ok := condition.Cdr.(*goscm.Pair)
			if !ok {
				return goscm.SCM_Nil, errors.New("Cond entry list too short")
			}

			if consequent_cell.Cdr != goscm.SCM_Nil {
				return goscm.SCM_Nil, errors.New("Trailing list elements in cond entry")
			}

			return consequent_cell.Car, nil
		}
		
		args, ok = args.Cdr.(*goscm.Pair)
		if !ok {
			return goscm.SCM_Nil, errors.New("Condition entries are a dotted list?")
		}
	}

	return goscm.SCM_Nil, nil
}

func scm_not(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Too many arguments")
	}
	
	if reflect.TypeOf(args.Car) == reflect.TypeOf(&goscm.Boolean{}) {
		return goscm.NewSCMT(!args.Car.(*goscm.Boolean).Value), nil
	}

	return goscm.NewSCMT(false), nil
}

func scm_reverse(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	if args == goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Expected exactly one argument")
	}
	if args.Cdr != goscm.SCM_Nil {
		return goscm.SCM_Nil, errors.New("Expected exactly one argument")
	}
	if reflect.TypeOf(args.Car) != reflect.TypeOf(&goscm.Pair{}) {
		return goscm.SCM_Nil, errors.New("Expected list type")
	}

	list, err := args.Car.(*goscm.Pair).ToSlice()
	if err != nil {
		return goscm.SCM_Nil, err
	}

	ret := goscm.SCM_Nil
	for i := 0; i < len(list); i++ {
		ret = goscm.Cons(list[i], ret)
	}
	return ret, nil
}

func scm_list(args *goscm.Pair, env *goscm.Environ) (goscm.SCMT, error) {
	return args, nil
}
