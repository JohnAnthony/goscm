package goscm

import (
	"strconv"
)

type SCMT_Integer struct {
	Value int
}

func (in *SCMT_Integer) String() string {
	return strconv.Itoa(in.Value)
}

func (in *SCMT_Integer) Eval(*SCMT_Env) SCMT {
	return in
}
