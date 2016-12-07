package user

import (
	"modules/user/aaa"

	"common/redis"
	"common/utils"
	"modules/user/config"

	"common/assert"

	"modules/misc/trans"

	"common/models/common"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type forgotPayload struct {
	Email string `json:"email" validate:"email"`
	usr   *aaa.User
}

func (fp *forgotPayload) Validate(ctx echo.Context) error {
	if err := validator.New().Struct(fp); err != nil {
		return err
	}

	m := aaa.NewAaaManager()
	u, err := m.FindUserByEmail(fp.Email)
	if err != nil {
		return trans.E("email not found")
	}

	fp.usr = u
	return nil
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
	key := <-utils.ID
	assert.Nil(aredis.StoreKey(key, pl.usr.Email, ucfg.Cfg.TokenTimeout))

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
	// TODO : change password to string when we are ready for it
	user.Password = common.NullString{String: pass, Valid: true}
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
