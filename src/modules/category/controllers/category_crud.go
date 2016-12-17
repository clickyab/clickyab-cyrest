package category

import (
	"common/assert"
	"modules/category/cat"
	"strings"

	"modules/misc/trans"

	"strconv"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type categoryPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" validate:"required"`
}

func (pl *categoryPayload) ValidateExtra(ctx echo.Context) error {
	if !cat.IsValidScope(pl.Scope) {
		return trans.E(`scope that is not valid in this app . you can this values "%s"`, strings.Join(cat.AllValidScopes(), `","`))
	}

	return nil
}

// createCategory
// @Route {
//		url	=	/create
//		method	=	post
//		payload	=	categoryPayload
//		resource=	manage_category:global
//		middleware = authz.Authenticate
//		200	=	cat.Category
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) createCategory(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*categoryPayload)
	m := cat.NewCatManager()
	c, err := m.Create(pl.Title, pl.Description, pl.Scope)
	assert.Nil(err)
	return u.OKResponse(ctx, c)
}

// editCategory
// @Route {
//		url	=	/edit/:id
//		method	=	put
//		payload	=	categoryPayload
//		resource=	manage_category:global
//		middleware = authz.Authenticate
//		200	=	cat.Category
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) editCategory(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*categoryPayload)
	id := ctx.Param("id")
	m := cat.NewCatManager()
	c, err := m.Update(pl.Title, pl.Description, pl.Scope, strconv.Atoi(id))
	assert.Nil(err)
	return u.OKResponse(ctx, c)
}
