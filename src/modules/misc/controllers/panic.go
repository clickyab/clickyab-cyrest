package misc

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// getVersion generate a panic for test
// @Route {
// 		url = /panic
//		method = get
//      500 = base.ErrorResponseSimple
//		200 = base.NormalResponse
// }
func (u *Controller) doPanic(ctx echo.Context) error {
	logrus.Panic("requested panic")
	return nil
}
