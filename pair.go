package goscm

type SCMT_Pair struct {
	car SCMT
	cdr SCMT
}

func (pair *SCMT_Pair) scm_eval(env *SCMT_Environment) *SCMT {
	// COMPLEX!
	return nil
}

func (pair *SCMT_Pair) scm_print() string {
	ret := "("
	ret += pair.car.scm_print()
	ret += " . "
	ret += pair.cdr.scm_print()
	ret += ")"

	return ret
}
