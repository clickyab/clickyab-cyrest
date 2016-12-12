package user

// AUTO GENERATED CODE DO NOT EDIT!

import (
	"common/assert"
	"common/controllers/base"
	"common/utils"
	"encoding/json"
	"modules/user/aaa"
	"modules/user/middlewares"
	"strings"

	"gopkg.in/labstack/echo.v3"
)

type listUserResponse struct {
	Total      int64                    `json:"total"`
	Data       []map[string]interface{} `json:"data"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	Definition base.Columns             `json:"definition"`
}

var listUserDefinition base.Columns

// @Route {
// 		url = /users
//		method = get
//		resource = user_list:parent
//		_sort_ = string, the sort and order like id:asc or id:desc available column "id","created_at","updated_at"
//		_user_type_ = string , filter the user_type field valid values are "personal","corpartion"
//		_status_ = string , filter the status field valid values are "registered","verified","blocked"
//		_email_ = string , search the email field
//      200 = listUserResponse
// }
func (u *Controller) listUser(ctx echo.Context) error {
	m := aaa.NewAaaManager()
	usr := authz.MustGetUser(ctx)
	p, c := utils.GetPageAndCount(ctx.Request(), false)

	filter := make(map[string]string)

	if e := ctx.Request().URL.Query().Get("user_type"); e != "" && aaa.UserType(e).IsValid() {
		filter["user_type"] = e
	}

	if e := ctx.Request().URL.Query().Get("status"); e != "" && aaa.UserStatus(e).IsValid() {
		filter["status"] = e
	}

	search := make(map[string]string)

	if e := ctx.Request().URL.Query().Get("email"); e != "" {
		search["email"] = e
	}

	s := ctx.Request().URL.Query().Get("sort")
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		parts = append(parts, "asc")
	}
	sort := parts[0]
	if !utils.StringInArray(sort, "id", "created_at", "updated_at") {
		sort = ""
	}
	order := strings.ToUpper(parts[1])
	if !utils.StringInArray(order, "ASC", "DESC") {
		order = "ASC"
	}

	pc := base.NewPermInterfaceComplete(usr, usr.ID, "user_list", "parent")
	dt, cnt := m.FillUserDataTableArray(pc, filter, search, sort, order, p, c)
	res := listUserResponse{
		Total:   cnt,
		Data:    dt.Filter(usr),
		Page:    p,
		PerPage: c,
	}
	if ctx.Request().URL.Query().Get("def") == "true" {
		res.Definition = listUserDefinition
	}
	return u.OKResponse(
		ctx,
		res,
	)
}

func init() {
	tmp := []byte(` [
		{
			"data": "parent_id",
			"name": "ParentID",
			"searchable": false,
			"sortable": false,
			"visible": false,
			"filter": false,
			"title": "ParentID",
			"filter_valid_map": null
		},
		{
			"data": "owner_id",
			"name": "OwnerID",
			"searchable": false,
			"sortable": false,
			"visible": false,
			"filter": false,
			"title": "OwnerID",
			"filter_valid_map": null
		},
		{
			"data": "id",
			"name": "ID",
			"searchable": false,
			"sortable": true,
			"visible": true,
			"filter": false,
			"title": "ID",
			"filter_valid_map": null
		},
		{
			"data": "email",
			"name": "Email",
			"searchable": true,
			"sortable": false,
			"visible": true,
			"filter": false,
			"title": "Email",
			"filter_valid_map": null
		},
		{
			"data": "user_type",
			"name": "Type",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": true,
			"title": "User type",
			"filter_valid_map": {
				"corpartion": "UserTypeCorporation",
				"personal": "UserTypePersonal"
			}
		},
		{
			"data": "avatar",
			"name": "Avatar",
			"searchable": false,
			"sortable": false,
			"visible": false,
			"filter": false,
			"title": "Avatar",
			"filter_valid_map": null
		},
		{
			"data": "status",
			"name": "Status",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": true,
			"title": "User status",
			"filter_valid_map": {
				"blocked": "UserStatusBlocked",
				"registered": "Registered",
				"verified": "UserStatusVerified"
			}
		},
		{
			"data": "created_at",
			"name": "CreatedAt",
			"searchable": false,
			"sortable": true,
			"visible": true,
			"filter": false,
			"title": "Created at",
			"filter_valid_map": null
		},
		{
			"data": "updated_at",
			"name": "UpdatedAt",
			"searchable": false,
			"sortable": true,
			"visible": true,
			"filter": false,
			"title": "Created at",
			"filter_valid_map": null
		}
	] `)
	assert.Nil(json.Unmarshal(tmp, &listUserDefinition))
}
