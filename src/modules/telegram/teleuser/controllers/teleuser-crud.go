package telegram

import (
	"common/assert"
	aredis "common/redis"
	"fmt"
	"math/rand"
	"modules/misc/trans"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"time"

	"strconv"

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

//	deleteTeleUser delete telegram user
//	@Route	{
//	url	=	/list/:id
//	method	= delete
//	resource = add_teleuser:self
//	middleware = authz.Authenticate
//	200 = base.NormalResponse
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) deleteTeleUser(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := tlu.NewTluManager()
	teleUser, err := m.FindTeleUserByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(teleUser.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("add_teleuser", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	m.DeleteTelegramUser(id)
	return u.OKResponse(ctx, nil)
}
