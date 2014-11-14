package goscm

import (
	"testing"
)

func Test_Bool(t *testing.T) {
	// Test true
	btrue := &Boolean { Value: true }
	if btrue.String() != "#t" {
		t.Error(btrue)
	}
	
	// Test false
	bfalse := &Boolean { Value: false }
	if bfalse.String() != "#f" {
		t.Error(bfalse)
	}
}
