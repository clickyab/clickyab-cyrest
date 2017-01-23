// +build ignore

package plan

import "modules/plan/pln"

type (
	plans []*pln.Plan
)

var _ = make(plans, 0) // make the type used :/
