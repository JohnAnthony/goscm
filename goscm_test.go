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
	pair := Cons(Make_SCMT(4), Make_SCMT(5))
	fmt.Println(pair.scm_print())
}
