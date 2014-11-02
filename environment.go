package goscm

type SCMT_Env struct {
	table map[string]SCMT
	parent *SCMT_Env
}

func (env *SCMT_Env) Eval(*SCMT_Env) SCMT {
	return env
}

func (*SCMT_Env) String() string {
	return "#<environment>"
}

func (env *SCMT_Env) Add(symb *SCMT_Symbol, val SCMT) {
	env.table[symb.value] = val
}

func (env *SCMT_Env) Find(symb *SCMT_Symbol) SCMT {
	ret := env.table[symb.value]
	if ret != nil {
		return ret
	}
	if env.parent == nil {
		return nil
	}
	return env.parent.Find(symb)
}

func (env *SCMT_Env) BindForeign(name string, f func (*SCMT_Pair, *SCMT_Env) SCMT) {
	env.Add(Make_Symbol(name), Make_Foreign(f))
}

func (env *SCMT_Env) BindSpecial(name string, f func (*SCMT_Pair, *SCMT_Env) SCMT) {
	env.Add(Make_Symbol(name), Make_Special(f))
}

// Environment provided #1: Completely empty. This is used for procedure
// environment closures.
// Procedures provided: NONE

func EnvEmpty(parent *SCMT_Env) *SCMT_Env {
	return &SCMT_Env {
		table: make(map[string]SCMT),
		parent: parent,
	}
}

// Environment provided #2: Simple. This is helpful for testing and as an
// example of how to build your own scheme env from scratch.
// Procedures provided: + - * / quote car cdr cons

func EnvSimple() *SCMT_Env {
	env := EnvEmpty(nil)
	
	env.BindForeign("+", scm_add)
	env.BindForeign("-", scm_subtract)
	env.BindForeign("*", scm_multiply)
	env.BindForeign("/", scm_divide)
	env.BindForeign("car", scm_car)
	env.BindForeign("cdr", scm_cdr)
	env.BindForeign("cons", scm_cons)
	env.BindForeign("map", scm_map)
	env.BindSpecial("quote", scm_quote)
	env.BindSpecial("define", scm_define)
	env.BindSpecial("begin", scm_begin)
	env.BindSpecial("let", scm_let)
	env.BindSpecial("lambda", scm_lambda)
	return env
}

// WARNING! These procedures do no input validation, so feeding them incorrect
// input will have strange effects!

func scm_add(args *SCMT_Pair, env *SCMT_Env) SCMT {
	ret := 0
	for ; !args.IsNil(); args = Cdr(args).(*SCMT_Pair) {
		ret += Car(args).(*SCMT_Integer).value
	}
	return Make_SCMT(ret)
}

func scm_multiply(args *SCMT_Pair, env *SCMT_Env) SCMT {
	ret := 1
	for ; !args.IsNil(); args = Cdr(args).(*SCMT_Pair) {
		ret *= Car(args).(*SCMT_Integer).value
	}
	return Make_SCMT(ret)
}

func scm_subtract(args *SCMT_Pair, env *SCMT_Env) SCMT {
	ret := Car(args).(*SCMT_Integer).value
	for args = Cdr(args).(*SCMT_Pair); !args.IsNil(); args = args.cdr.(*SCMT_Pair) {
		ret -= Car(args).(*SCMT_Integer).value
	}
	return Make_SCMT(ret)
}

func scm_divide(args *SCMT_Pair, env *SCMT_Env) SCMT {
	ret := Car(args).(*SCMT_Integer).value
	for args = Cdr(args).(*SCMT_Pair); !args.IsNil(); args = args.cdr.(*SCMT_Pair) {
		ret /= Car(args).(*SCMT_Integer).value
	}
	return Make_SCMT(ret)
}

func scm_car(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return args.car.(*SCMT_Pair).car
}

func scm_cdr(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return args.car.(*SCMT_Pair).cdr
}

func scm_cons(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return Cons(args.car, args.cdr.(*SCMT_Pair).car)
}

func scm_map(args *SCMT_Pair, env *SCMT_Env) SCMT {
	ret := SCMT_Nil
	f := Car(args).Eval(env).(SCMT_Func)
	target := Car(Cdr(args).(*SCMT_Pair)).(*SCMT_Pair)
	for l := target; l != SCMT_Nil; l = Cdr(l).(*SCMT_Pair) {
		arg := Cons(Car(l), SCMT_Nil)
		ret = Cons(f.Apply(arg, env), ret)
	}
	return Reverse(ret)
}

func scm_quote(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return args.car
}

func scm_define(args *SCMT_Pair, env *SCMT_Env) SCMT {
	env.Add(Car(args).(*SCMT_Symbol), Car(Cdr(args).(*SCMT_Pair)).Eval(env))
	return SCMT_Nil 
}

func scm_begin(args *SCMT_Pair, env *SCMT_Env) SCMT {
	var result SCMT
	for result = SCMT_Nil; args != SCMT_Nil; args = Cdr(args).(*SCMT_Pair) {
		result = Car(args).Eval(env)
	}
	return result
}

func scm_let(args *SCMT_Pair, env *SCMT_Env) SCMT {
	body := Cdr(args).(*SCMT_Pair)
	newenv := EnvEmpty(env)
	
	for vars := Car(args).(*SCMT_Pair); vars != SCMT_Nil; vars = Cdr(vars).(*SCMT_Pair) {
		symb := Car(Car(vars).(*SCMT_Pair)).(*SCMT_Symbol)
		val := Car(Cdr(Car(vars).(*SCMT_Pair)).(*SCMT_Pair)).Eval(env)
		newenv.Add(symb, val)
	}

	return scm_begin(body, newenv)
}

func scm_lambda(args *SCMT_Pair, env *SCMT_Env) SCMT {
	return Make_Proc(Car(args).(*SCMT_Pair), Cdr(args).(*SCMT_Pair), env)
}
