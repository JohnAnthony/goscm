package goscm

import (
	"testing"
)

func Test_Read(t *testing.T) {
	// Numeric atom
	inst := goscm.NewInstance()
	v := inst.Read("2")
	//if (v.Type != SCM_Number)
	// error
	//if (v.Value != 2)
	// error
}
