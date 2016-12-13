package aaa


// Domain model
// @Model {
//		table = domains
//		primary = true, id
//		find_by = id,cname
//		transaction = insert
//		list = yes
// }
type Domain struct {
	ID int64 `json:"id" db:"id"`
	CName string `json:"cname" db:"cname"`
}

