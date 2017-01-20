package ads

import "common/models/common"

// Plan model
// @Model {
//		table = plans
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Plan struct {
	ID          int64             `db:"id" json:"id" sort:"true" title:"ID"`
	Name        string            `json:"name" db:"name" title:"Name"`
	Description common.NullString `json:"description" db:"description" title:"Description"`
}
