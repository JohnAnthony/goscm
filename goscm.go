package goscm

import (
)

type SCMT interface {
	scm_eval(*SCMT_Environment) SCMT
	scm_print() string
}
