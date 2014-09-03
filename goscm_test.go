package goscm

import (
	"testing"
)

func Test_Read(t *testing.T) {
	inst := NewInstance("")
	var ret *SCMType
	var remain string
	
	// Numbers //

	ret, remain = inst.Read("2")
	if ret.Type != SCM_Number {
		t.Error("Expected to be of type SCM_Number")
	} else if *ret.Value.(*int) != 2 {
    	t.Error("Expected returned value to be 2")
	} else if remain != "" {
    	t.Error("Expected no remainder")
	}

	// ret, remain = inst.Read("2 3 4 5")
	// if (ret.Type != SCM_Number)
	// error
	// if (ret.Value != 2)
	// error
	// if (remain != " 3 4 5")
	// error
	
	// Pairs //

	// ret, remain = inst.Read("(1 2 3)")
	// if (ret.Type != SCM_Pair)
	// error
	// if (ret.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Car.Value != 1)
	// error
	// if (ret.Value.Cdr.Type != SCM_Pair)
	// error
	// if (ret.Value.Cdr.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Cdr.Value.Car.Value != 2)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Type != SCM_Pair)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Car.Value != 3)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Cdr != nil)
	// error
}

func Test_Eval(t *testing.T) {
}

func Test_Print(t *testing.T) {
}
