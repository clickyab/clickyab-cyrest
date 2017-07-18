package user

import (
	"common/assert"
	"common/mail"
	"common/redis"
	"common/utils"
	"modules/misc/trans"
	"modules/user/aaa"

	"common/config"

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

	// TODO we need this three lines later for better recovery workflow
	//key := <-utils.ID
	//assert.Nil(aredis.StoreKey(key, usr.Email, ucfg.Cfg.TokenTimeout))
	//sendEmailCodeGen(usr,key)

	pass := utils.PasswordGenerate(8)
	usr.Password = pass
	assert.Nil(m.UpdateUser(usr))
	sendEmailPasswordGen(usr, pass)

	return u.OKResponse(
		ctx,
		nil,
	)
}

//func sendEmailCodeGen(usr *aaa.User, key string) {
//	link := fmt.Sprintf("%s://%s/v1/new-password/%s", config.Config.Proto, config.Config.Site, key)
//
//	mail.SendByTemplateName(trans.T("Password recovery").Translate("fa_IR"), "recoverCode", struct {
//		Name string
//		Link string
//	}{
//		usr.Email,
//		link,
//	}, config.Config.Mail.From, usr.Email)
//}

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
	sendEmailPasswordGen(user, pass)
	return u.OKResponse(
		ctx,
		nil,
	)
}

func sendEmailPasswordGen(usr *aaa.User, pass string) {
	err := mail.SendByTemplateName(trans.T("Password recovery").Translate("fa_IR"), "forgot-new-password", struct {
		Name string
		Pass string
	}{
		usr.Email,
		pass,
	}, config.Config.Mail.From, usr.Email)
	assert.Nil(err)
}
