package middlewares

import (
	"encoding/json"
	"modules/misc/trans"
	"net/http"
	"reflect"

	"common/assert"
	"common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ContextBody is the context key for the body unmarshalled object
	ContextBody string = "_body"
)

// Validator is used to validate the input
type Validator interface {
	// Validate return true, nil on no error, false ,error map in error string
	Validate(*gin.Context) (bool, map[string]string)
}

// PayloadUnMarshallerGenerator create a middleware base on the pattern for the request body
func PayloadUnMarshallerGenerator(pattern interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Make sure the request is POST or PUT since DELETE and GET must not have payloads
		method := strings.ToUpper(c.Request.Method)
		assert.True(
			!utils.StringInArray(method, "GET", "DELETE"),
			"[BUG] Get and Delete must not have request body",
		)
		// Create a copy
		cp := reflect.New(reflect.TypeOf(pattern)).Elem().Addr().Interface()
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(cp)
		if err != nil {
			c.Header("error", trans.T("invalid request body"))
			e := struct {
				Error string `json:"error"`
			}{
				Error: "invalid request body",
			}

			c.JSON(http.StatusBadRequest, e)
			c.Abort()
			return
		}
		if valid, ok := cp.(Validator); ok {
			if ok, errs := valid.Validate(c); ok {
				c.Set(ContextBody, cp)
			} else {
				c.Header("error", trans.T("invalid request body"))
				c.JSON(http.StatusBadRequest, errs)
				c.Abort()
				return
			}
		} else { // Just add it, no validation
			c.Set(ContextBody, cp)
		}
		c.Next()
	}
}

// GetPayload from the request
func GetPayload(c *gin.Context) (interface{}, bool) {
	return c.Get(ContextBody)
}
