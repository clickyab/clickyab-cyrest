package misc

import (
	"modules/misc/t9n"

	"modules/misc/trans"

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
	return u.OKResponse(ctx, translation)

}
