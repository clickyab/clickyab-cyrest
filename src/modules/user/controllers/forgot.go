package user

import (
	"modules/user/aaa"

	"common/redis"
	"common/utils"
	"modules/user/config"

	"common/assert"

	"github.com/labstack/echo"
)

type forgotPayload struct {
	Email string `json:"email"`
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
	//check exist email
	pl := u.MustGetPayload(ctx).(*forgotPayload)
	m := aaa.NewAaaManager()
	_, err := m.FindUserByEmail(pl.Email)
	if err != nil {
		return u.BadResponse(ctx, err)
	}
	//generate key for email
	key := <-utils.ID
	//insert code:email in redis
	assert.Nil(aredis.StoreKey(key, pl.Email, ucfg.Cfg.TokenTimeout))
	//send email
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
	m := aaa.NewAaaManager()
	user, err := m.FindUserByEmail(email)
	if err != nil {
		return u.BadResponse(ctx, err)
	}
	pass := utils.PasswordGenerate(8)
	//change password
	user.Password.String = pass
	//change token
	user.AccessToken = <-utils.ID
	m.UpdateUser(user)
	sendEmailPasswordGen()
	return u.OKResponse(
		ctx,
		nil,
	)
}

func sendEmailPasswordGen() {
	//todo send email to user
}
