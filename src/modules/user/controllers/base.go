package user

import (
	"common/assert"
	"common/controllers/base"
	"modules/user/aaa"
	"modules/user/middlewares"

	"github.com/labstack/echo"
)

type BaseController struct {
	base.Controller
}

// GetUser try to get back the user to system
func (c BaseController) MustGetUser(ctx echo.Context) *aaa.User {
	u, ok := authz.GetUser(ctx)
	assert.True(ok, "[BUG] user is not in the context")

	return u
}

func init() {
	base.Register(&Controller{})
}
