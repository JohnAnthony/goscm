package goscm

// External package for type/integer
// External package for type/string
// External package for pair traversal

import (
)

type SCMT interface {
	scm_eval(*SCMT_Environment) SCMT
	scm_print() string
}
