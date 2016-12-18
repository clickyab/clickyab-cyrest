package controllers

import (
	"modules/location/loc"

	"gopkg.in/labstack/echo.v3"
)

// listCountry
// @Route {
//		url	=	/country
//		method	=	get
//		middleware = authz.Authenticate
//		200	=	loc.Country
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) listCountry(ctx echo.Context) error {
	m := loc.NewLocManager()
	country := m.ListCountries()
	return u.OKResponse(ctx, country)

}
