package goscm

import (
	"reflect"
)

func Make_SCMT(in interface {}) SCMT {
	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return &SCMT_Integer { value: in.(int) }
	case reflect.String:
		return &SCMT_String { value: in.(string) }
	default:
		return nil
	}
}
