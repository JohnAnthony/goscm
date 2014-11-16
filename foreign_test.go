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
		t.Error()
	}

	// Check that it returns the correct retuls
	sq, err := scm_f.Apply(NewList(NewPlainInt(13)), EnvEmpty(nil))
	if err != nil {	t.Error(err) }
	ans, ok := sq.(*PlainInt)
	if !ok { t.Error(ok) }
	if ans.String() != "169" {
		t.Error(ans)
	}
}
