package base

import (
	"common/assert"
	"common/middlewares"
	"net/http"

	"common/try"

	"github.com/gin-gonic/gin"
)

// NormalResponse is for 2X responses
type NormalResponse struct {
}

// ComplexResponse for the result, when the result type in not in the structure
type ComplexResponse map[string]interface{}

// ErrorResponseMap is the map for the response with detail error mapping
type ErrorResponseMap map[string]string

// ErrorResponseSimple is the type for response when the error is simply a string
type ErrorResponseSimple struct {
	Error string `json:"error"`
}

// Controller is the base controller for all other controllers
type Controller struct {
}

// BadResponse is 400 request
func (c Controller) BadResponse(ctx *gin.Context, err error) {
	err = try.Try(err)
	ctx.Header("error", err.Error())
	ctx.JSON(http.StatusBadRequest, ErrorResponseSimple{Error: err.Error()})
}

// NotFoundResponse is 404 request
func (c Controller) NotFoundResponse(ctx *gin.Context, err error) {
	var res = ErrorResponseSimple{}
	if err != nil {
		res.Error = try.Try(err).Error()
	} else {
		res.Error = http.StatusText(http.StatusNotFound)
	}
	ctx.Header("error", res.Error)
	ctx.JSON(http.StatusNotFound, res)
}

// OKResponse is 200 request
func (c Controller) OKResponse(ctx *gin.Context, res interface{}) {
	if res == nil {
		res = NormalResponse{}
	}
	ctx.JSON(http.StatusOK, res)
}

// MustGetPayload is for payload middleware
func (c Controller) MustGetPayload(ctx *gin.Context) interface{} {
	obj, ok := middlewares.GetPayload(ctx)
	assert.True(ok, "[BUG] payload un-marshaller failed")

	return obj
}
