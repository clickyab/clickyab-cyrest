package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"common/mail"

	"time"

	"gopkg.in/labstack/echo.v3"
)

//@Validate{
// }
type registrationPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gt=5"`
}

// registerUser register user in system
// @Route {
// 		url = /register
//		method = post
//      payload = registrationPayload
//		200 = responseLoginOK
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) registerUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*registrationPayload)
	m := aaa.NewAaaManager()

	usr, err := m.RegisterUser(pl.Email, pl.Password)
	if err != nil {
		return u.BadResponse(ctx, trans.E("email is already registered in our system"))
	}

	token := m.GetNewToken(usr, ctx.Request().UserAgent(), ctx.RealIP())

	mail.SendByTemplateName(trans.T("شما با موفقیت در روبیک اد ثبت شدید").String(), "resource/register.html", struct {
		Date time.Time
		Name string
	}{
		Date: time.Now(),
		Name: pl.Email,
	}, "info@rubikad.com", "guest@sdfs.com")

	return u.OKResponse(
		ctx,
		createLoginResponse(usr, token),
	)
}
