package goscm

import (
	"reflect"
	"errors"
)

func EnsureAll(pair *SCMT_Pair, t reflect.Type) error {
	for pair != SCMT_Nil {
		if reflect.TypeOf(pair.Car) != t {
			return errors.New("Wrong type")
		}
		if reflect.TypeOf(pair.Cdr) != reflect.TypeOf(&SCMT_Pair{}) {
			return errors.New("Dotted list")
		}
		pair = pair.Cdr.(*SCMT_Pair)
	}

	return nil
}
