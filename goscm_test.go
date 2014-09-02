package goscm

import (
	"testing"
)

func Test_Read(t *testing.T) {
	inst := NewInstance(nil)
	var ret *SCMCell
	var remain string
	
	// Numbers //

	ret, remain = inst.Read("2")
	//if (ret.Type != SCM_Number)
	// error
	// if (ret.Value != 2)
	// error
	// if (remain != "")
	// error

	ret, remain = inst.Read("2 3 4 5")
	// if (ret.Type != SCM_Number)
	// error
	// if (ret.Value != 2)
	// error
	// if (remain != " 3 4 5")
	// error
	
	// Lists //

	ret, remain = inst.Read("(1 2 3)")
	// if (ret.Type != SCM_List)
	// error
	// if (ret.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Car.Value != 1)
	// error
	// if (ret.Value.Cdr.Type != SCM_List)
	// error
	// if (ret.Value.Cdr.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Cdr.Value.Car.Value != 2)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Type != SCM_List)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Car.Type != SCM_Number)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Car.Value != 3)
	// error
	// if (ret.Value.Cdr.Value.Cdr.Value.Cdr != nil)
	// error
}
