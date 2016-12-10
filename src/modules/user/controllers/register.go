package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"gopkg.in/labstack/echo.v3"
)

type personal struct {
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Cellphone string            `json:"cellphone"`
	Gender    aaa.ProfileGender `json:"gender"`
}

type corporation struct {
	CompanyName string `json:"company_name"`
	Phone       string `json:"phone"`
}

type registrationPayload struct {
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Source      aaa.UserSource `json:"source"`
	Personal    *personal      `json:"personal" validate:"-"`
	Corporation *corporation   `json:"corporation" validate:"-"`
}

func (r *registrationPayload) Validate(ctx echo.Context) error {
	if r.Personal != nil && r.Corporation != nil {
		return trans.E("both personal and corporation is set")
	}
	if r.Personal == nil && r.Corporation == nil {
		return trans.E("both personal and corporation is not set")
	}

	//validator.New().Struct()
	return nil
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

	var prof interface{}
	if pl.Personal != nil {
		prof = aaa.NewUserProfilePersonal(pl.Personal.FirstName, pl.Personal.LastName, pl.Personal.Gender, pl.Personal.Cellphone)
	} else {
		prof = aaa.NewUserProfileCorporation(pl.Corporation.CompanyName, pl.Corporation.Phone)
	}

	user, err := m.RegisterUser(pl.Email, pl.Password, prof)
	if err != nil {
		return u.BadResponse(ctx, err)
	}

	token := m.GetNewToken(user.AccessToken)
	return u.OKResponse(
		ctx,
		responseLoginOK{
			UserID:      user.ID,
			Email:       user.Email,
			AccessToken: token,
		},
	)
}
