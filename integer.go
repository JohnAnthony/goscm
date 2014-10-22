package goscm

import (
	"strconv"
)

type SCMT_Integer struct {
	value int
}

func (in *SCMT_Integer) String() string {
	return strconv.Itoa(in.value)
}

func (in *SCMT_Integer) scm_eval(*SCMT_Env) SCMT {
	return in
}
