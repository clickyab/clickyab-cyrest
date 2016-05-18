package user

import (
	"common/assert"
	"fmt"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/templates"
	"modules/user/utils"
	"modules/user/utils/mailer"
	"modules/user/utils/sms"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type reserveUserPayload struct {
	Contact string `json:"contact"`
}

type reservedTokenResponse struct {
	//	Token    string `json:"token"`
	OldToken bool `json:"old_token"`
}

// reserveUser is the route for reserve a email/phone for registration
// @Route {
// 		url = /challenge
//		method = post
//		payload = reserveUserPayload
//      200 = reservedTokenResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) challengeCreate(ctx *gin.Context) {
	pl := u.MustGetPayload(ctx).(*reserveUserPayload)
	m := aaa.NewAaaManager()
	token, old, err := m.ReserveToken(pl.Contact)
	if err != nil {
		u.BadResponse(ctx, err)
		return
	}
	// its time to send the message to user
	// TODO : better function for handling this transparent
	t, err := utils.DetectContactType(pl.Contact)
	if t == utils.TypeEmail {
		var body string
		body, err = templates.RenderMail(
			"base",
			struct {
				Subject string
				Message string
			}{
				Subject: trans.T("Registration"),
				Message: fmt.Sprintf(trans.T("Your code is %s"), token.Token),
			},
		)
		assert.Nil(err)
		go func() {
			err = mailer.SendMail(pl.Contact, pl.Contact, body, trans.T("Registration confirm"))
			if err != nil {
				logrus.Error(err)
			}
		}()
	} else if t == utils.TypePhone {
		go func() {
			err = sms.SendSMS(pl.Contact, fmt.Sprintf("Your code is %s ", token.Token))
			if err != nil {
				logrus.Error(err)
			}
		}()
	}

	u.OKResponse(
		ctx,
		reservedTokenResponse{
			//Token:    token.Token,
			OldToken: old,
		},
	)
}
