package goscm

import "testing"

func Test_Pair(t *testing.T) {
	// Printing the singleton (556677)
	sing := NewList(NewPlainInt(556677))
	if sing.String() != "(556677)" {
		t.Error(sing)
	}
	
	// Printing of the pair (4 . 5)
	pair := NewList(NewPlainInt(4), NewPlainInt(5))
	if pair.String() != "(4 . 5)" {
		t.Error(pair)
	}
	
	// A test with the list (1 2 3 4 5 6 7 8 9) aka
	// (1 . (2 . (3 . (4 . (5 . (6 . (7 . (8 . (9 . ())))))))))
	list := NewList (
		NewPlainInt(9),
		NewPlainInt(8),
		NewPlainInt(7),
		NewPlainInt(6),
		NewPlainInt(5),
		NewPlainInt(4),
		NewPlainInt(3),
		NewPlainInt(2),
		NewPlainInt(1),
	)
	if list.String() != "(1 2 3 4 5 6 7 8 9)" {
		t.Error(list)
	}
}

func Test_Nil(t *testing.T) {
	n := SCM_Nil
	
	// Test how it prints
	if n.String() != "()" {
		t.Error(n)
	}
}
