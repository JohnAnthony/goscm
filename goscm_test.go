package goscm

import (
	"testing"
)

func Test_Read(t *testing.T) {
	// Numeric atom
	inst := NewInstance()
	var ret *SCMCell
	var remain string
	
	ret, remain = inst.Read("2")
	//if (ret.Type != SCM_Number)
	// error
	//if (ret.Value != 2)
	// error

	ret, remain = inst.Read("2 3 4 5")
	//if (ret.Type != SCM_Number)
	// error
	//if (ret.Value != 2)
	// error
	//if (remain != "3 4 5")
	// error
}
