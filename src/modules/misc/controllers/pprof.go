package misc

import (
	"runtime/pprof"

	"gopkg.in/labstack/echo.v3"
)

// getVersion generate a panic for test
// @Route {
// 		url = /pprof
//		method = get
//      500 = base.ErrorResponseSimple
//		200 = base.NormalResponse
// }
func (u *Controller) doPProf(ctx echo.Context) error {
	return pprof.Lookup("goroutine").WriteTo(ctx.Response(), 1)
}
