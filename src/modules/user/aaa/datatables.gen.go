package aaa

import (
	"common/controllers/base"
	"strings"
)

// AUTO GENERATED CODE DO NOT EDIT!

type (
	UserDataTableArray []UserDataTable
)

func (udta UserDataTableArray) Filter(u base.PermInterface) []map[string]interface{} {
	res := make([]map[string]interface{}, len(udta))
	for i := range udta {
		res[i] = udta[i].Filter(u)
	}

	return res
}

// Filter is for filtering base on permission
func (udt UserDataTable) Filter(u base.PermInterface) map[string]interface{} {
	res := map[string]interface{}{

		"parent_id": udt.ParentID,

		"owner_id": udt.OwnerID,

		"id": udt.ID,

		"email": udt.Email,

		"source": udt.Source,

		"user_type": udt.Type,

		"avatar": udt.Avatar,

		"status": udt.FormatStatus(),

		"created_at": udt.CreatedAt,

		"updated_at": udt.UpdatedAt,
	}

	action := []string{}

	if _, ok := u.HasPermStringOn("user_edit", udt.OwnerID, udt.ParentID, "global"); ok {
		action = append(action, "edit")
	}

	res["_actions"] = strings.Join(action, ",")
	return res
}
