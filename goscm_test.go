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
    	t.Error("Expected returned value to be 2")
	} else if remain != "" {
    	t.Error("Expected no remainder")
	}

	// ret, remain = inst.Read("2 3 4 5")
	// if (ret.Type != SCM_Integer)
	// error
	// if (ret.Value != 2)
	// error
	// if (remain != " 3 4 5")
	// error
}
