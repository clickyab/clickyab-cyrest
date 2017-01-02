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
//		url	=	/
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
	c := m.Create(pl.Title, pl.Description, pl.Scope)
	return u.OKResponse(ctx, c)
}

// editCategory
// @Route {
//		url	=	/:id
//		method	=	put
//		payload	=	categoryPayload
//		resource=	manage_category:global
//		middleware = authz.Authenticate
//		200	=	cat.Category
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) editCategory(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*categoryPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	m := cat.NewCatManager()
	c, err := m.FindCategoryByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	c.Title = pl.Title
	c.Description = pl.Description
	c.Scope = pl.Scope

	assert.Nil(m.UpdateCategory(c))

	return u.OKResponse(ctx, c)
}

//	getCategory
//	@Route	{
//	url	=	/:id
//	method	= get
//	resource = list_Category:global
//	middleware = authz.Authenticate
//	200 = cat.Category
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getCategory(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := cat.NewCatManager()
	category, err := m.FindCategoryByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	return u.OKResponse(ctx, category)
}
