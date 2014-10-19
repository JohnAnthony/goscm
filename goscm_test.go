package goscm

import (
	"testing"
	"fmt"
)

func Test_Integer(t *testing.T) {
	in := SCMT_Integer { value: 10 }
	fmt.Println(in.scm_print())
}

func Test_String(t *testing.T) {
	str := SCMT_String { value: "foobar!" }
	fmt.Println(str.scm_print())
}

func Test_Pair(t *testing.T) {
	pair := Cons(&SCMT_Integer { value: 4 }, &SCMT_Integer { value: 5 })
	fmt.Println(pair.scm_print())
}
