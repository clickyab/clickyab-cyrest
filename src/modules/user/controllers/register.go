package user

import (
	"modules/user/aaa"

	"github.com/labstack/echo"
	"regexp"
)

type responseLoginOK struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	AccessToken string `json:"Accesstoken"`
}

type registrationPayload struct {
	Email    string `json:"email" validate:"min=3,max=40"`
	Password string `json:"password"`
	Personal bool   `json:"personal"`

	FirstName   string `json:"first_name" validate:"nonzero"`
	LastName    string `json:"last_name" validate:"nonzero"`
	CompanyName string `json:"company_name"`

	Cellphone string `json:"cellphone" validate:"nonzero"`
	Phone     string `json:"phone"`
	Gender    int    `json:"gender"`
	Source    string `json:"source"`
}

func (r *registrationPayload) Validate(ctx echo.Context) (bool, map[string]string) {
	var res = make(map[string]string)
	var fail bool

	if len(r.Password) < 6 {
		res["password"] = "Password is invalid"
		fail = true
	}
	valid,_ := ValidateEmail(r.Email)
	if len(r.Email) < 4 || !valid {
		res["email"] = "Email is invalid"
		fail = true
	}
	if len(r.Phone) < 6 {
		res["phone"] = "Phone is invalid"
		fail = true
	}

	if r.Personal {
		if len(r.FirstName) < 2 {
			res["firstname"] = "First name is invalid"
			fail = true
		}
		if len(r.LastName) < 2 {
			res["lastName"] = "Last name is invalid"
			fail = true
		}
		if len(r.Cellphone) < 2 {
			res["cellphone"] = "Cellphone is invalid"
			fail = true
		}

	} else {
		if len(r.CompanyName) < 2 {
			res["companyName"] = "Company Name is invalid"
			fail = true
		}

	}

	if fail {
		return false, res
	}

	return true, nil
}
func ValidateEmail(email string) (bool,error) {
	//Re := regexp.MustCompile()
	return regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)
}

// registerUser register user in system
// @Route {
// 		url = /register
//		method = post
//      payload = registrationPayload
//      200 = responseLoginOK
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) registerUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*registrationPayload)
	m := aaa.NewAaaManager()

	user, err := m.RegisterUser(pl.Email, pl.Password, pl.Personal)
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
