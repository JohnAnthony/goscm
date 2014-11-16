package goscm

import "testing"

func Test_Special(t *testing.T) {
	env := EnvEmpty(nil)
	
	//
	// A non-expanding special form
	//
	quote_spesh := NewSpecial(func (args *Pair, env *Environ) (SCMT, error) {
		return args, nil
	})
	env.Add(NewSymbol("quote"), quote_spesh)

	list := NewList(
		NewSymbol("quote"),
		NewSymbol("a"),
		NewSymbol("b"),
		NewSymbol("c"),
		NewSymbol("d"),
		NewSymbol("e"),
	)
	
	ret, err := list.Eval(env)
	if err != nil {	t.Error(err) }
	if ret.String() != "(A B C D E)" { t.Error(ret)	}
	
	//
	// An expanding special form
	//
	scm_if := func (args *Pair, env *Environ) (SCMT, error) {
		argss, err := args.ToSlice()
		if err != nil { return SCM_Nil, err }

		pred, err := argss[0].Eval(env)
		if err != nil { return SCM_Nil, err }

		if pred.(*Boolean).Value {
			return argss[1], nil
		}

		return argss[2], nil
	}
	
	if_spesh := NewSpecialTCO(scm_if)
	env.Add(NewSymbol("if"), if_spesh)
	env.Add(NewSymbol("a"), NewPlainInt(22234))
	env.Add(NewSymbol("b"), NewPlainInt(33345))

	list2 := NewList(
		NewSymbol("if"),
		NewBoolean(true),
		NewSymbol("a"),
		NewSymbol("b"),
	)
	
	// Check it applies properly
	exeval, err := list2.Eval(env)
	if err != nil {	t.Error(err) }
	if exeval.String() != "22234" { t.Error(exeval) }

	// Check it expands properly
	exex, err := if_spesh.Expand(list2.Cdr.(*Pair), env)
	if err != nil { t.Error(err) }
	if exex.String() != "A" { t.Error(exex) }
}
