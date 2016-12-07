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
	"gopkg.in/go-playground/validator.v9"
)

const (
	// ContextBody is the context key for the body unmarshalled object
	ContextBody string = "_body"
)

// Validator is used to validate the input
type Validator interface {
	// Validate return error if the type is invalid
	Validate(echo.Context) error
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
					Error: "invalid request body",
				}

				c.JSON(http.StatusBadRequest, e)
				return err
			}
			if valid, ok := cp.(Validator); ok {
				if errs := valid.Validate(c); errs == nil {
					c.Set(ContextBody, cp)
				} else {
					c.Request().Header().Set("error", trans.T("invalid request body"))
					if ve, ok := errs.(validator.ValidationErrors); ok {
						tmp := make(map[string]string)
						for i := range ve {
							tmp[ve[i].Field()] = ve[i].Translate(nil)
						}
						return c.JSON(http.StatusBadRequest, errs)
					}
					return c.JSON(http.StatusBadRequest, errs)
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
