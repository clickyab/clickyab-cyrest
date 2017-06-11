package user

import (
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"
	"time"

	"gopkg.in/labstack/echo.v3"
)

type profilePayload struct {
	Profile *profile `json:"profile"`
}

// @Validate {
// }
type profile struct {
	FirstName    string            `json:"first_name" validate:"gt=2" error:"first name must be valid"`
	LastName     string            `json:"last_name" validate:"gt=2" error:"last name must be valid"`
	Birthday     time.Time         `json:"birthday"`
	Gender       aaa.ProfileGender `json:"gender"`
	CellPhone    string            `json:"cellphone"`
	Phone        string            `json:"phone"`
	Address      string            `json:"address"`
	ZipCode      string            `json:"zip_code"`
	NationalCode string            `json:"national_code"`
	CountryID    int64             `json:"country_id"`
	ProvinceID   int64             `json:"province_id"`
	CityID       int64             `json:"city_id"`
}

// editProfile
// @Route {
//		url	=	/profile
//		method	=	post
//		payload	= profilePayload
//		middleware = authz.Authenticate
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) editProfile(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*profilePayload)
	m := aaa.NewAaaManager()
	user := authz.MustGetUser(ctx)
	token := authz.MustGetToken(ctx)
	if pl.Profile != nil { //the profile has been selected
		_, err := m.RegisterProfile(
			user.ID,
			pl.Profile.FirstName,
			pl.Profile.LastName,
			pl.Profile.Birthday,
			pl.Profile.Gender,
			pl.Profile.CellPhone,
			pl.Profile.Phone,
			pl.Profile.Address,
			pl.Profile.ZipCode,
			pl.Profile.NationalCode,
			pl.Profile.CountryID,
			pl.Profile.ProvinceID,
			pl.Profile.CityID,
		)
		if err != nil {
			return u.BadResponse(ctx, trans.E("can not update profile"))
		}

		return u.OKResponse(
			ctx,
			createLoginResponse(user, token),
		)
	}
	return u.BadResponse(ctx, trans.E("can not update profile"))
}
