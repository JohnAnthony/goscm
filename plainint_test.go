package goscm

import "testing"

func Test_PlainInt(t *testing.T) {
	in := Make_PlainInt(31337)
	if in.Value != 31337 {
		t.Error(in.Value)
	}
	if in.String() != "31337" {
		t.Error(in)
	}
}
