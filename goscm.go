package goscm

import "reflect"

type SCMT interface {
	scm_eval(*SCMT_Environment) SCMT
	String() string
}

func Make_SCMT(in interface {}) SCMT {
	if in == nil {
		return &SCMT_Nil {}
	}

	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return &SCMT_Integer { value: in.(int) }
	case reflect.String:
		return &SCMT_String { value: in.(string) }
	default:
		// TODO: We probably need to put an error here
		return nil
	}
}
