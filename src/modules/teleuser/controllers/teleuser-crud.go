package telegram

import (
	"common/assert"
	aredis "common/redis"
	"fmt"
	"math/rand"
	"modules/teleuser/tlu"
	"modules/user/middlewares"
	"time"

	"gopkg.in/labstack/echo.v3"
)

//	add add teleuser
//	@Route	{
//		url = /
//		method = post
//		resource = add_teleuser:self
//		200 = tlu.Verifycode
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) add(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	num := rand.Intn(90000000) + 10000000

	key := fmt.Sprintf("%d:%d", currentUser.ID, num)
	assert.Nil(aredis.StoreKey(key, fmt.Sprintf("%d", currentUser.ID), 48*time.Hour))
	return u.OKResponse(ctx, tlu.Verifycode{Key: key})
}
