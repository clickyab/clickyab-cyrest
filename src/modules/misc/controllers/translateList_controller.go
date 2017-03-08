package misc

import (
	"common/assert"
	"common/utils"
	"encoding/json"
	"modules/misc/base"
	"modules/misc/t9n"

	"modules/misc/trans"

	"gopkg.in/labstack/echo.v3"
)

type listTranslatelistResponse struct {
	Total      int64                       `json:"total"`
	Data       t9n.TranslateDataTableArray `json:"data"`
	Page       int                         `json:"page"`
	PerPage    int                         `json:"per_page"`
	Definition base.Columns                `json:"definition"`
}

var listTranslatelistDefinition base.Columns

// @Route {
// 		url = /translate/:lang
//		method = get
//		resource = translate_list:global
//		_lang_ = string, the language
//		_def_ = bool, show definition in result?
//		_c_ = int , count per page
//		_p_ = int , page number
//		200 = listTranslatelistResponse
// }
func (u *Controller) listTranslatelist(ctx echo.Context) error {
	m := t9n.NewT9nManager()
	p, c := utils.GetPageAndCount(ctx.Request(), false)

	filter := make(map[string]string)

	search := make(map[string]string)
	lang := ctx.Param("lang")
	if lang != trans.PersianLang && lang != trans.EnglishLang {
		return u.NotFoundResponse(ctx, trans.E("not valid language"))
	}
	sort := ""
	order := "ASC"

	//pc := base.NewPermInterfaceComplete(usr, usr.ID, "translate_list", "global")
	dt, cnt := m.FillTranslateDataTableArray(lang, filter, search, sort, order, p, c)
	res := listTranslatelistResponse{
		Total:   cnt,
		Data:    dt,
		Page:    p,
		PerPage: c,
	}
	if ctx.Request().URL.Query().Get("def") == "true" {
		res.Definition = listTranslatelistDefinition
	}
	return u.OKResponse(
		ctx,
		res,
	)
}


// @Route {
// 		url = /milad
//		method = get
//		_lang_ = string, the language
//		_def_ = bool, show definition in result?
//		_c_ = int , count per page
//		_p_ = int , page number
//		200 = listTranslatelistResponse
// }
func (u *Controller) miladTranslatelist(ctx echo.Context) error {
	panic("paniced")
	m := t9n.NewT9nManager()
	p, c := utils.GetPageAndCount(ctx.Request(), false)

	filter := make(map[string]string)

	search := make(map[string]string)
	lang := ctx.Param("lang")
	if lang != trans.PersianLang && lang != trans.EnglishLang {
		return u.NotFoundResponse(ctx, trans.E("not valid language"))
	}
	sort := ""
	order := "ASC"

	//pc := base.NewPermInterfaceComplete(usr, usr.ID, "translate_list", "global")
	dt, cnt := m.FillTranslateDataTableArray(lang, filter, search, sort, order, p, c)
	res := listTranslatelistResponse{
		Total:   cnt,
		Data:    dt,
		Page:    p,
		PerPage: c,
	}
	if ctx.Request().URL.Query().Get("def") == "true" {
		res.Definition = listTranslatelistDefinition
	}
	return u.OKResponse(
		ctx,
		res,
	)
}

func init() {
	tmp := []byte(` [
		{
			"data": "id",
			"name": "ID",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": false,
			"title": "ID",
			"filter_valid_map": null
		},
		{
			"data": "text",
			"name": "Text",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": false,
			"title": "Text",
			"filter_valid_map": null
		},
		{
			"data": "lang",
			"name": "Lang",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": false,
			"title": "Lang",
			"filter_valid_map": null
		},
		{
			"data": "translated",
			"name": "Translated",
			"searchable": false,
			"sortable": false,
			"visible": true,
			"filter": false,
			"title": "Translated",
			"filter_valid_map": null
		},
		{
			"data": "_actions",
			"name": "Actions",
			"searchable": false,
			"sortable": false,
			"visible": false,
			"filter": false,
			"title": "Actions",
			"filter_valid_map": null
		}
	] `)
	assert.Nil(json.Unmarshal(tmp, &listTranslatelistDefinition))
}
