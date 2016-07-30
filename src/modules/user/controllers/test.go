package user

import "github.com/labstack/echo"

type (
	tmp struct {
		ID   int64  `json:"id"`
		Text string `json:"text"`
		userListResponse
		Data []string `json:"data"`
	}

	tmpArray []tmp
)

// testFunction test code generator functionality
// @Route {
// 		url = /test
//		method = post
//      payload = tmp
//		#resource = resource_name
//      200 = tmpArray
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) testFunction(ctx echo.Context) error {
	return nil
}
