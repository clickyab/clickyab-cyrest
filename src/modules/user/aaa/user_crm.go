package aaa

// UserCRM model
// @Model {
//		table = user_crm
//		primary = false, user_id
//		find_by = user_id
// }
type UserCRM struct {
	UserID          int64  `db:"user_id" json:"user_id"`
	OriginatingLead string `db:"originating_lead" json:"originating_lead"`
	CustomerCode    string `db:"originating_lead" json:"originating_lead"`
	GID             string `db:"gid" json:"gid"`
	LeadStatus      int    `db:"lead_status" json:"lead_status"`
	ReadByCRM       int    `db:"read_by_crm" json:"read_by_crm"`
}
