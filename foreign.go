package goscm

type SCMT_Foreign struct {
	function func (*SCMT_Pair) SCMT
}

func (fo *SCMT_Foreign) scm_eval(*SCMT_Env) SCMT {
	return fo
}

func (*SCMT_Foreign) String() string {
	return "#<foreign function>"
}

func (fo *SCMT_Foreign) Apply(args *SCMT_Pair) SCMT {
	return fo.function(args) // TODO
}

func Make_Foreign(f func (*SCMT_Pair) SCMT) *SCMT_Foreign {
	return &SCMT_Foreign {
		function: f,
	}
}
