package goscm

type SCMT_Environment struct {
	table map[string]SCMT
	child *SCMT_Environment
}

func (*SCMT_Environment) scm_eval(*SCMT_Environment) SCMT {
	// TODO: This shoudl probably be an error
	return nil
}

func (*SCMT_Environment) String() string {
	return "#<environment>"
}

func EnvEmpty(child *SCMT_Environment) *SCMT_Environment {
	return &SCMT_Environment {
		table: make(map[string]SCMT),
		child: child,
	}
}
