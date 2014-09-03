package goscm

import (
)

type Instance struct {
}

type SCMType struct {
	Type int
	Value interface {}
}

const (
	SCM_Error = iota
	SCM_Number
)

// NewInstance takes a string description of a standard and returns a new scheme
// interpreter instance initialized to meet that standard.
// If the string description is "", then we return a very minimal instance.
func NewInstance(std string) *Instance {
	return &Instance {
	}
}

// Read takes a string and reads it using the instance's syntax rules.
// It returns a *SCMType of the first sexp it encounters and a string containing
// the unprocessed remainder of the input string.
// If EOF is reached while attempting to read an sexp, then an appropriate
// scheme error is returned
func (inst *Instance) Read(s string) (*SCMType, string) {
	return &SCMType{ Type: SCM_Error, Value: "Not implemented" }, ""
}

// Eval takes a *SCMType and evaluates it inside of an environment
func (inst *Instance) Eval(*SCMType) *SCMType {
	return nil
}

// Print takes a *SCMType and returns a string representation of that value
// Note the this function does not actually print the value to the screen
func (inst *Instance) Print(*SCMType) string {
	return ""
}
