package goscm

import "testing"

func Test_Pair(t *testing.T) {
	// Printing the singleton (556677)
	sing := Make_List(Make_PlainInt(556677))
	if sing.String() != "(556677)" {
		t.Error(sing)
	}
	
	// Printing of the pair (4 . 5)
	pair := Make_List(Make_PlainInt(4), Make_PlainInt(5))
	if pair.String() != "(4 . 5)" {
		t.Error(pair)
	}
	
	// A test with the list (1 2 3 4 5 6 7 8 9) aka
	// (1 . (2 . (3 . (4 . (5 . (6 . (7 . (8 . (9 . ())))))))))
	list := Make_List (
		Make_PlainInt(9),
		Make_PlainInt(8),
		Make_PlainInt(7),
		Make_PlainInt(6),
		Make_PlainInt(5),
		Make_PlainInt(4),
		Make_PlainInt(3),
		Make_PlainInt(2),
		Make_PlainInt(1),
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
