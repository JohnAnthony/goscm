package goscm

import (
	"testing"
)

func Test_Read_Integer(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var remain string
	
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
		t.Error("Expected remainder to be \" 30 40 50\", got", remain)
	}
}

func Test_Read_String(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var remain string

	inst = NewInstance("")
	ret, remain = inst.Read("\"Test string\"")
	if ret.Type != SCM_String {
		t.Error("Expected to be of type SCM_String")
	} else if *ret.Value.(*string) != "Test string" {
    	t.Error("Expected returned value to be \"Test string\", got", *ret.Value.(*string))
	} else if remain != "" {
    	t.Error("Expected no remainder, got", remain)
	}
	
	inst = NewInstance("")
	ret, remain = inst.Read("\"Test string\" 2 3 4")
	if ret.Type != SCM_String {
		t.Error("Expected to be of type SCM_String")
	} else if *ret.Value.(*string) != "Test string" {
    	t.Error("Expected returned value to be \"Test string\", got", *ret.Value.(*string))
	} else if remain != " 2 3 4" {
    	t.Error("Expected remainder to be \" 2 3 4\", got", remain)
	}
}
