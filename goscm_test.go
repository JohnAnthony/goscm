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
	}
	if *ret.Value.(*int) != 2 {
    	t.Error("Expected returned value to be 2, got", *ret.Value.(*int))
	}
	if remain != "" {
    	t.Error("Expected no remainder, got", remain)
	}

	ret, remain = inst.Read("90 30 40 50")
	if ret.Type != SCM_Integer {
		t.Error("Expected to be of type SCM_Integer")
	}
	if *ret.Value.(*int) != 90 {
    	t.Error("Expected returned value to be 90, got", *ret.Value.(*int))
	} 
	if remain != " 30 40 50" {
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
	}
	if *ret.Value.(*string) != "Test string" {
    	t.Error("Expected returned value to be \"Test string\", got", *ret.Value.(*string))
	}
	if remain != "" {
    	t.Error("Expected no remainder, got", remain)
	}
	
	inst = NewInstance("")
	ret, remain = inst.Read("\"Test string\" 2 3 4")
	if ret.Type != SCM_String {
		t.Error("Expected to be of type SCM_String")
	}
	if *ret.Value.(*string) != "Test string" {
    	t.Error("Expected returned value to be \"Test string\", got", *ret.Value.(*string))
	}
	if remain != " 2 3 4" {
    	t.Error("Expected remainder to be \" 2 3 4\", got", remain)
	}

	inst = NewInstance("")
	ret, remain = inst.Read("\"\"")
	if ret.Type != SCM_String {
		t.Error("Expected to be of type SCM_String")
	}
	if *ret.Value.(*string) != "" {
    	t.Error("Expected returned value to be \"\", got", *ret.Value.(*string))
	}
	if remain != "" {
    	t.Error("Expected remainder to be \"\", got", remain)
	}

	// TODO: Test for handling of unterminated strings
}

func Test_Read_Symbol(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var remain string

	inst = NewInstance("")
	ret, remain = inst.Read("foo")
	if ret.Type != SCM_Symbol {
		t.Error("Expected to be of type SCM_Symbol")
	}
	if *ret.Value.(*string) != "FOO" {
		t.Error("Expected returned value to be \"FOO\", got", *ret.Value.(*string))
	}
	if remain != "" {
		t.Error("Expected no remainder, got", remain)
	}

	inst = NewInstance("")
	ret, remain = inst.Read("abrax ebran ubrol irwin")
	if ret.Type != SCM_Symbol {
		t.Error("Expected to be of type SCM_Symbol")
	}
	if *ret.Value.(*string) != "ABRAX" {
		t.Error("Expected returned value to be \"ABRAX\", got", *ret.Value.(*string))
	}
	if remain != " ebran ubrol irwin" {
		t.Error("Expected no remainder, got", remain)
	}
}

func Test_Read_List(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var remain string
	
	inst = NewInstance("")
	ret, remain = inst.Read("(69)")
	if ret.Type != SCM_Pair {
		t.Error("Expected to be of type SCM_Pair")
	}
	if ret.Value.(*SCMPair).Car.Type != SCM_Integer {
		t.Error("Expected list's car to be a SCMInteger")
	}
	if *ret.Value.(*SCMPair).Car.Value.(*int) != 69 {
		t.Error("Expected list's car to be 69, got",
		*ret.Value.(*SCMPair).Car.Value.(*int))
	}
	if ret.Value.(*SCMPair).Cdr != nil {
		t.Error("The list's cdr is not nil")
	}
	if remain != "" {
		t.Error("Expected no remainder, got", remain)
	}

	inst = NewInstance("")
	ret, remain = inst.Read("(+ 10 20)")
	head := ret
	if head.Type != SCM_Pair {
		t.Error("Expected to be of type SCM_Pair")
	}
	if head.Value.(*SCMPair).Car.Type != SCM_Symbol {
		t.Error("Expected list's car to be a SCMPair")
	}
	if *head.Value.(*SCMPair).Car.Value.(*string) != "+" {
		t.Error("Expected list's car to be symbol \"+\"")
	}
	if head.Value.(*SCMPair).Cdr.Type != SCM_Pair {
		t.Error("Expected list's cdr to be a SCM_Pair")
	}
	head = head.Value.(*SCMPair).Cdr
	if head.Value.(*SCMPair).Car.Type != SCM_Integer {
		t.Error("Expected list's cdar to be a SCM_Integer")
	}
	if *head.Value.(*SCMPair).Car.Value.(*int) != 10 {
		t.Error(
			"Expected list's cdar value to be 10, got",
			*head.Value.(*SCMPair).Car.Value.(*int))
	}
	if head.Value.(*SCMPair).Cdr.Type != SCM_Pair {
		t.Error("Expected list's next element to be a SCM_Pair")
	}
	head = head.Value.(*SCMPair).Cdr
	if head.Value.(*SCMPair).Car.Type != SCM_Integer {
		t.Error("Expected list's cdar to be a SCM_Integer")
	}
	if *head.Value.(*SCMPair).Car.Value.(*int) != 20 {
		t.Error(
			"Expected list's cddar value to be 20, got",
			*head.Value.(*SCMPair).Car.Value.(*int))
	}
	if head.Value.(*SCMPair).Cdr != nil {
		t.Error("Expected next list element to be nil")
	}
	if remain != "" {
		t.Error("Expected no remainder, got", remain)
	}
	
	inst = NewInstance("")
	ret, remain = inst.Read("()")
	if ret != nil {
		t.Error("Expected return value to be nil (an empty list), got", ret)
	}
	if remain != "" {
		t.Error("Expected no remainder, got", remain)
	}
	
	// TODO: Add test for remainders
}

func Test_Print_Integer(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var p string

	inst = NewInstance("")
	ret, _ = inst.Read("10")
	p = inst.Print(ret)
	if p != "10" {
		t.Error("Expected Read(\"10\") to return 10, got", p)
	}

	inst = NewInstance("")
	ret, _ = inst.Read("9984523")
	p = inst.Print(ret)
	if p != "9984523" {
		t.Error("Expected Read(\"9984523\") to return 9984523, got", p)
	}
}

func Test_Print_String(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var p string

	inst = NewInstance("")
	ret, _ = inst.Read("\"foo\"")
	p = inst.Print(ret)
	if p != "\"foo\"" {
		t.Error("Expected string \"foo\" to evaluate to itself in quotes, got", p)
	}
}

func Test_Print_Symbol(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var p string

	inst = NewInstance("")
	ret, _ = inst.Read("foo")
	p = inst.Print(ret)
	if p != "FOO" {
		t.Error("Expected symbol foo to evaluate to FOO, got", p)
	}
}

func Test_Print_List(t *testing.T) {
	var inst *Instance
	var ret *SCMType
	var p string

	inst = NewInstance("")
	ret, _ = inst.Read("(+ 10 20 30 40)")
	p = inst.Print(ret)
	if p != "(+ 10 20 30 40)" {
		t.Error("Expected (+ 10 20 30 40), got", p)
	}
	
	inst = NewInstance("")
	ret, _ = inst.Read("(foo (sub list) 45 5)")
	p = inst.Print(ret)
	if p != "(FOO (SUB LIST) 45 5)" {
		t.Error("Expected (FOO (SUB LIST) 45 5), got", p)
	}
}
