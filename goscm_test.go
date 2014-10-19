package goscm

import (
	"testing"
	"fmt"
)

func Test_Integer(t *testing.T) {
	in := Make_SCMT(10)
	fmt.Println(in.scm_print())
}

func Test_String(t *testing.T) {
	str := Make_SCMT("foobar!")
	fmt.Println(str.scm_print())
}

func Test_Pair(t *testing.T) {
	pair := Cons(Make_SCMT("Foobarrrr!"), Make_SCMT(nil))
	fmt.Println(pair.scm_print())
	
	pair2 := Cons(Make_SCMT(4), Make_SCMT(5))
	fmt.Println(pair2.scm_print())
	
	list := Cons(Make_SCMT(10), Make_SCMT(nil))
	list = Cons(Make_SCMT(9), list)
	list = Cons(Make_SCMT(8), list)
	list = Cons(Make_SCMT(7), list)
	list = Cons(Make_SCMT(6), list)
	list = Cons(Make_SCMT(5), list)
	list = Cons(Make_SCMT(4), list)
	list = Cons(Make_SCMT(3), list)
	list = Cons(Make_SCMT(2), list)
	list = Cons(Make_SCMT(1), list)
	fmt.Println(list.scm_print())
}
