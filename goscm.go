package goscm

import (
	"strconv"
)

type Instance struct {
}

type SCMType struct {
	Type SCMT
	Value interface {}
}

type SCMT int

const (
	SCM_Integer SCMT = iota
	SCM_String
	SCM_Symbol
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
	var ret SCMType
	var start int
	var end int

	for start = 0; s[start] == ' '; start++ {}
	
	if (s[start] >= '0' && s[start] <= '9') { // A number
		ret.Type = SCM_Integer
		var val int
		for end = start; end < len(s) && s[end] != ' '; end++ {}
		val, _ = strconv.Atoi(s[start:end])
		ret.Value = &val
	} else if (s[start] == '"') { // A string
		ret.Type = SCM_String
		var val string
		start++
		for end = start; end < len(s) && s[end] != '"'; end++ {}
		val = s[start:end]
		end++
		ret.Value = &val
	} else { // A symbol
		ret.Type = SCM_Symbol
		var val string
		for end = start; end < len(s) && s[end] != ' '; end++ {}
		val = s[start:end]
		ret.Value = &val
	}

	return &ret, s[end:]
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
