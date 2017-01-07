package location

import (
	"modules/location/loc"
	"strconv"

	"fmt"
	"gopkg.in/labstack/echo.v3"
)

// listProvince
// @Route {
//		url	=	/province
//		method	=	get
//		middleware = authz.Authenticate
//		200	=	loc.Province
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) listProvince(ctx echo.Context) error {
	m := loc.NewLocManager()
	province := m.ListProvinces()
	return u.OKResponse(ctx, province)
}

// listProvinceByCountry
// @Route {
//		url	=	/province/:id
//		method	=	get
//		middleware = authz.Authenticate
//		200	=	loc.Province
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) listProvinceByCountry(ctx echo.Context) error {
	m := loc.NewLocManager()
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	fmt.Println(id)
	province := m.ListProvinceByCountryID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	return u.OKResponse(ctx, province)
}
