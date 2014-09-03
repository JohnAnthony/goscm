package goscm

import (
	"testing"
)

func Test_Read_Integer(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var remain string
	
	// Integers //

	inst = NewInstance("")
	ret, remain = inst.Read("2")
	if ret.Type != SCM_Integer {
		t.Error("Expected to be of type SCM_Integer")
	} else if *ret.Value.(*int) != 2 {
    	t.Error("Expected returned value to be 2, got", *ret.Value.(*int))
	} else if remain != "" {
    	t.Error("Expected no remainder, got", remain)
	}

	ret, remain = inst.Read("90 30 40 50")
	if ret.Type != SCM_Integer {
		t.Error("Expected to be of type SCM_Integer")
	} else if *ret.Value.(*int) != 90 {
    	t.Error("Expected returned value to be 90, got", *ret.Value.(*int))
	} else if remain != " 30 40 50" {
		t.Error("Expected remainder to be \" 30 40 50\", got:", remain)
	}
}
