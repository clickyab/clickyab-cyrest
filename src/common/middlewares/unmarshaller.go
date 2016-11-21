package middlewares

import (
	"encoding/json"
	"modules/misc/trans"
	"net/http"
	"reflect"

	"common/assert"
	"common/utils"
	"strings"

	"github.com/labstack/echo"
)

const (
	// ContextBody is the context key for the body unmarshalled object
	ContextBody string = "_body"
)

// Validator is used to validate the input
type Validator interface {
	// Validate return true, nil on no error, false ,error map in error string
	Validate(echo.Context) (bool, map[string]string)
}

// PayloadUnMarshallerGenerator create a middleware base on the pattern for the request body
func PayloadUnMarshallerGenerator(pattern interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Make sure the request is POST or PUT since DELETE and GET must not have payloads
			method := strings.ToUpper(c.Request().Method())
			assert.True(
				!utils.StringInArray(method, "GET", "DELETE"),
				"[BUG] Get and Delete must not have request body",
			)
			// Create a copy
			cp := reflect.New(reflect.TypeOf(pattern)).Elem().Addr().Interface()
			decoder := json.NewDecoder(c.Request().Body())
			err := decoder.Decode(cp)
			if err != nil {
				c.Request().Header().Set("error", trans.T("invalid request body"))
				e := struct {
					Error string `json:"error"`
				}{
					Error: "invalid request body1",
				}

				c.JSON(http.StatusBadRequest, e)
				return err
			}
			if valid, ok := cp.(Validator); ok {
				if ok, errs := valid.Validate(c); ok {
					c.Set(ContextBody, cp)
				} else {
					c.Request().Header().Set("error", trans.T("invalid request body2"))
					c.JSON(http.StatusBadRequest, errs)
					return trans.E("invalid request body3")
				}
			} else {
				// Just add it, no validation
				c.Set(ContextBody, cp)
			}
			return next(c)
		}
	}
}

// GetPayload from the request
func GetPayload(c echo.Context) (interface{}, bool) {
	t := c.Get(ContextBody)
	return t, t != nil
}
