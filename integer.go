package goscm

import (
	"strconv"
)

type SCMT_Integer struct {
	value int
}

func (in *SCMT_Integer) scm_print() string {
	return strconv.Itoa(in.value)
}

func (in *SCMT_Integer) scm_eval(*SCMT_Environment) SCMT {
	return in
}
