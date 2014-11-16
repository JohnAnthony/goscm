package goscm

import (
	"strconv"
)

type PlainInt struct {
	Value int
}

func (in *PlainInt) String() string {
	return strconv.Itoa(in.Value)
}

func (in *PlainInt) Eval(*Environ) (SCMT, error) {
	return in, nil
}

func Make_PlainInt(n int) *PlainInt {
	return &PlainInt { Value: n }
}