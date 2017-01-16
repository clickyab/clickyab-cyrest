package user

import (
	"common/assert"
	"common/redis"
	"common/utils"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/config"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type forgotPayload struct {
	Email string `json:"email" validate:"email"`
}

// forgotPassword get email
// @Route {
//		url	=	/forgot/call
//		method	=	post
//		payload	= forgotPayload
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) forgotPassword(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*forgotPayload)
	//generate key for email

	m := aaa.NewAaaManager()
	usr, err := m.FindUserByEmail(pl.Email)
	if err != nil {
		return u.BadResponse(ctx, trans.E("email not found"))
	}

	key := <-utils.ID
	assert.Nil(aredis.StoreKey(key, usr.Email, ucfg.Cfg.TokenTimeout))

	sendEmailCodeGen()
	return u.OKResponse(
		ctx,
		nil,
	)
}

func sendEmailCodeGen() {
	//todo send email to user
}

// forgotGeneratePassword get email
// @Route {
//		url	=	/forgot/callback/:code
//		method	=	get
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) forgotGeneratePassword(ctx echo.Context) error {
	code := ctx.Param("code")
	email, err := aredis.GetKey(code, false, 0)
	assert.Nil(err)
	m := aaa.NewAaaManager()
	user, err := m.FindUserByEmail(email)
	if err != nil {
		return u.BadResponse(ctx, trans.E("email is not registered"))
	}
	pass := utils.PasswordGenerate(8)
	// TODO : change password to string when we are ready for it
	user.Password = pass
	assert.Nil(m.UpdateUser(user))
	sendEmailPasswordGen()
	return u.OKResponse(
		ctx,
		nil,
	)
}

func sendEmailPasswordGen() {
	//todo send email to user
}
