package teleuser

import (
	"common/redis"
	"fmt"
	"math/rand"
	"modules/user/middlewares"
	"time"

	"modules/teleuser/tlu"

	"gopkg.in/labstack/echo.v3"
)

//	add add teleuser
//	@Route	{
//	url	=	/
//	method	= post
//	resource = add_teleuser:self
//	middleware = authz.Authenticate
//	200 = tlu.Verifycode
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) add(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	num := rand.Intn(90000000)
	if num <= 9999999 {
		num = num + 10000000
	}
	key := fmt.Sprintf("%d:%d", currentUser.ID, num)
	aredis.StoreKey(key, fmt.Sprintf("%d", currentUser.ID), 48*time.Hour)
	return u.OKResponse(ctx, tlu.Verifycode{Key: key})

}
