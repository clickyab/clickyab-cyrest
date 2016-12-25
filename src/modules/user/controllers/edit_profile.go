package user

import (
	"modules/user/aaa"
	"time"

	"modules/user/middlewares"

	"modules/misc/trans"

	"modules/misc/middlewares"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/labstack/echo.v3"
)

type ProfilePayload struct {
	Personal    *Personal    `json:"personal"`
	Corporation *Corporation `json:"corporation"`
}

// Validate custom validation for editing profile
func (lp *ProfilePayload) Validate(ctx echo.Context) error {
	if (lp.Personal == nil && lp.Corporation == nil) || (lp.Personal != nil && lp.Corporation != nil) {
		return middlewares.GroupError{
			"personal":    trans.E("invalid payload body"),
			"corporation": trans.E("invalid payload body"),
		}
	}
	return validator.New().Struct(lp)
}

// @Validate {
// }
type Personal struct {
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

// @Validate {
// }
type Corporation struct {
	Title        string `json:"title" validate:"gt=3" error:"title must be valid"`
	EconomicCode string `json:"economic_code"`
	RegisterCode string `json:"register_code"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	CountryID    int64  `json:"country_id"`
	ProvinceID   int64  `json:"province_id"`
	CityID       int64  `json:"city_id"`
}

// editProfile
// @Route {
//		url	=	/profile
//		method	=	post
//		payload	= ProfilePayload
//		middleware = authz.Authenticate
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) editProfile(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*ProfilePayload)
	m := aaa.NewAaaManager()
	user := authz.MustGetUser(ctx)
	token := authz.MustGetToken(ctx)
	if pl.Personal != nil { //the personal profile has been selected
		_, err := m.RegisterPersonal(
			user.ID,
			pl.Personal.FirstName,
			pl.Personal.LastName,
			pl.Personal.Birthday,
			pl.Personal.Gender,
			pl.Personal.CellPhone,
			pl.Personal.Phone,
			pl.Personal.Address,
			pl.Personal.ZipCode,
			pl.Personal.NationalCode,
			pl.Personal.CountryID,
			pl.Personal.ProvinceID,
			pl.Personal.CityID,
		)
		if err != nil {
			return u.BadResponse(ctx, trans.E("can not update profile"))
		}

		return u.OKResponse(
			ctx,
			createLoginResponse(user, token),
		)
		return nil
	}

	if pl.Corporation != nil { ////the corporation profile has been selected
		_, err := m.RegisterCorporation(
			user.ID,
			pl.Corporation.Title,
			pl.Corporation.EconomicCode,
			pl.Corporation.RegisterCode,
			pl.Corporation.Phone,
			pl.Corporation.Address,
			pl.Corporation.CountryID,
			pl.Corporation.ProvinceID,
			pl.Corporation.CityID,
		)
		if err != nil {
			return u.BadResponse(ctx, trans.E("can not update profile"))
		}

		return u.OKResponse(
			ctx,
			createLoginResponse(user, token),
		)
		return nil
	}

	return nil
}
