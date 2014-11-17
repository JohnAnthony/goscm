package goscm

import "testing"

func Test_Foreign(t *testing.T) {
	f := func (list *Pair, env *Environ) (SCMT, error) {
		n := list.Car.(*PlainInt).Value
		return NewSCMT(n * n), nil
	}
	scm_f := NewForeign(f)
	
	// Check that it prints prettily
	if scm_f.String() != "#<foreign function>" {
		t.Error(scm_f)
	}

	// Check that it returns the correct retuls
	sq, err := scm_f.Apply(NewList(NewPlainInt(13)), NewEnv(nil))
	if err != nil {	t.Error(err) }
	ans, ok := sq.(*PlainInt)
	if !ok { t.Error(ok) }
	if ans.String() != "169" {
		t.Error(ans)
	}
}

func Test_Foreign_List(t *testing.T) {
	list := NewList(
		NewSymbol("+"),
		NewSCMT(190),
		NewSCMT(3),
		NewSCMT(11),
	)

	env := NewEnv(nil)
	add_func := func (args *Pair, env *Environ) (SCMT, error) {
		ret := 0
		for ; args != SCM_Nil; args = args.Cdr.(*Pair) {
			ret += args.Car.(*PlainInt).Value
		}
		return NewSCMT(ret), nil
	}
	env.Add(NewSymbol("+"), NewForeign(add_func))

	ret, err := list.Eval(env)
	if err != nil { t.Error(err) }
	if ret.String() != "204" { t.Error(ret) }
}
