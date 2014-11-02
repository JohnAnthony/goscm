package goscm

import "reflect"

type SCMT interface {
	Eval(*SCMT_Env) SCMT
	String() string
}

type SCMT_Func interface {
	Apply(*SCMT_Pair, *SCMT_Env) SCMT
}

func Make_SCMT(in interface {}) SCMT {
	if in == nil {
		return SCMT_Nil
	}

	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return &SCMT_Integer { Value: in.(int) }
	case reflect.String:
		return &SCMT_String { Value: in.(string) }
	default:
		// TODO: We probably need to put an error here
		return nil
	}
}
