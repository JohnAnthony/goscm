package goscm

import "testing"

func Test_Symbol(t *testing.T) {
	s := NewSymbol("Foo-bar")
	if s.String() != "FOO-BAR" {
		t.Error(s)
	}
}
