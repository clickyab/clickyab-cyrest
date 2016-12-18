package location

import (
	"modules/location/loc"
	"strconv"

	"gopkg.in/labstack/echo.v3"
)

// listCity
// @Route {
//		url	=	/city
//		method	=	get
//		middleware = authz.Authenticate
//		200	=	loc.City
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) listCity(ctx echo.Context) error {
	m := loc.NewLocManager()
	city := m.ListCities()
	return u.OKResponse(ctx, city)
}

// listCityByProvince
// @Route {
//		url	=	/city/:id
//		method	=	get
//		middleware = authz.Authenticate
//		200	=	loc.City
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) listCityByProvince(ctx echo.Context) error {
	m := loc.NewLocManager()
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	city := m.ListCityByProvinceID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	return u.OKResponse(ctx, city)
}
