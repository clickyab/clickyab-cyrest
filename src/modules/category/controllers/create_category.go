package category

import (
	"common/assert"
	cat "modules/category/cat"
	"modules/misc/trans"

	"fmt"
	"strings"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type createCategoryPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" validate:"required"`
}

// createCategory
// @Route {
//		url	=	/create-category
//		method	=	post
//		payload	=	createCategoryPayload
//		resource=	create-category
//		middleware = authz.Authenticate
//		200	=	cat.Category
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) createCategory(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*createCategoryPayload)
	if !cat.IsValidScope(pl.Scope) {
		s := fmt.Sprintf(`scope that is not valid in this app . you can this values "%s"`, strings.Join(cat.AllValidScopes(), `","`))
		return u.BadResponse(ctx, trans.E(s))
	}
	c := &cat.Category{Title: pl.Title, Description: pl.Description, Scope: pl.Scope}
	m := cat.NewCatManager()
	assert.Nil(m.CreateCategory(c))
	return u.OKResponse(ctx, c)

}
