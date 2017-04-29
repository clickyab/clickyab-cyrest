package misc

import (
	"modules/misc/t9n"

	"modules/misc/trans"

	"common/rabbit"
	"modules/telegram/bot/worker"

	"strconv"

	"common/assert"

	"gopkg.in/labstack/echo.v3"
)

type transPayload struct {
	StringID   int64  `json:"id" validate:"required"`
	Lang       string `json:"lang" validate:"required"`
	Translated string `json:"translated" validate:"required"`
}

//	translate
//	@Route	{
//	url	=	/translate
//	method	= put
//	payload	= transPayload
//	resource = trans_message:global
//	middleware = authz.Authenticate
//	200 = t9n.Translations
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) translate(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*transPayload)
	if pl.Lang != trans.PersianLang && pl.Lang != trans.EnglishLang {
		return u.NotFoundResponse(ctx, trans.E("language not supported"))
	}
	m := t9n.NewT9nManager()
	translation := &t9n.Translations{
		Lang:       pl.Lang,
		Translated: pl.Translated,
		StringID:   pl.StringID,
	}
	err := m.CreateOnDuplicateUpdateTranslations(translation)
	if err != nil {
		return u.NotFoundResponse(ctx, trans.E("error while inserting translation"))
	}
	trans.Clear()
	rabbit.MustPublish(&bot.ClearTrans{})
	return u.OKResponse(ctx, translation)

}

//	delete
//	@Route	{
//	url	=	/delete/:id
//	method	= delete
//	resource = delete_trans:global
//	middleware = authz.Authenticate
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) delete(ctx echo.Context) error {
	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 0)
	if err != nil {
		u.NotFoundResponse(ctx, err)
	}

	m := t9n.NewT9nManager()
	err = m.DeleteStringByID(id)
	assert.Nil(err)

	return u.OKResponse(ctx, "Deleted successfully")

}

//	transDump
//	@Route	{
//	url	=	/dump/:lang
//	method	= get
//	200 = t9n.Mixed
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) transDump(ctx echo.Context) error {
	lang := ctx.Param("lang")
	if lang != trans.PersianLang && lang != trans.EnglishLang {
		return u.NotFoundResponse(ctx, trans.E("language not supported"))
	}

	result := t9n.NewT9nManager().DumpAll(lang)

	return u.OKResponse(ctx, result)

}

type getTranslatePayload struct {
	Translate string `json:"translate" validate:"required"`
}

//	callIdentifyAd call identify
//	@Route	{
//		url	=	/translate
//		method	= post
//		payload	= getTranslatePayload
//		200 = base.NormalResponse
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) getTranslate(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*getTranslatePayload)
	trans.T(pl.Translate)
	return u.OKResponse(ctx, nil)
}
