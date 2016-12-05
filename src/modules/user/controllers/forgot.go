package user

import "github.com/labstack/echo"

type forgotPayload struct {
	Email string `json:"email"`
}
type responseForgotOK struct {
	message string `json:"string"`
}

//	forgotPassword get email
// 	@Route {
//		url	=	/forgot
//		method	=	post
//		payload	= forgotPayload
//		200	=	responseForgotOK
//		400	=	base.ErrorResponseSimple
//	}

func (u *Controller) forgotPassword(ctx echo.Context) error {
	return nil
}
