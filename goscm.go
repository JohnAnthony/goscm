package goscm

import (
	"strconv"
	"strings"
)

type Instance struct {
}

type SCMType struct {
	Type SCMT
	Value interface {}
}

type SCMPair struct {
	Car *SCMType
	Cdr *SCMType
}

type SCMT int

const (
	SCM_Integer SCMT = iota
	SCM_String
	SCM_Symbol
	SCM_Pair
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
	var end int

	s = strings.TrimLeft(s, " ")
	
	switch {
	case s[0] == '(': // A list
		tip := &ret
		s = s[1:]
		remain := strings.TrimLeft(s, " ")
		if remain[0] == ')' {
			return nil, remain[1:]
		}
		for {
			tip.Type = SCM_Pair
			tip.Value = new(SCMPair)
			tip.Value.(*SCMPair).Car, remain = inst.Read(remain)
			remain = strings.TrimLeft(remain, " ")
			if remain[0] == ')' {
				return &ret, remain[1:]
			}
			tip.Value.(*SCMPair).Cdr = new(SCMType)
			tip = tip.Value.(*SCMPair).Cdr
		}
	case s[0] >= '0' && s[0] <= '9': // A number
		ret.Type = SCM_Integer
		var val int
		for end = 0; end < len(s) && s[end] != ' ' && s[end] != ')'; end++ {}
		val, _ = strconv.Atoi(s[:end])
		ret.Value = &val
	case s[0] == '"': // A string
		ret.Type = SCM_String
		var val string
		s = s[1:]
		for end = 0; end < len(s) && s[end] != '"'; end++ {}
		val = s[:end]
		end++
		ret.Value = &val
	default: // A symbol
		ret.Type = SCM_Symbol
		var val string
		for end = 0; end < len(s) && s[end] != ' ' && s[end] != ')'; end++ {}
		val = strings.ToUpper(s[:end])
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
func (inst *Instance) Print(cell *SCMType) string {
	switch cell.Type {
	case SCM_Integer:
		return strconv.Itoa(*cell.Value.(*int))
	case SCM_String:
		return "\"" + *cell.Value.(*string) + "\""
	case SCM_Symbol:
		return *cell.Value.(*string)
	case SCM_Pair:
		ret := "("
		for ; cell != nil; {
			ret += inst.Print(cell.Value.(*SCMPair).Car)
			if cell.Value.(*SCMPair).Cdr != nil {
				ret += " "
			}
			cell = cell.Value.(*SCMPair).Cdr
		}
		return ret + ")"
	}
	
	return "ERROR UNPRINTABLE"
}
