package base

import (
	"common/assert"
	"common/middlewares"
	"errors"
	"net/http"

	"common/try"

	"gopkg.in/labstack/echo.v3"
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
func (c Controller) BadResponse(ctx echo.Context, err error) error {
	err = try.Try(err)
	ctx.Response().Header().Add("error", err.Error())
	ctx.JSON(http.StatusBadRequest, ErrorResponseSimple{Error: err.Error()})

	return err
}

// NotFoundResponse is 404 request
func (c Controller) NotFoundResponse(ctx echo.Context, err error) error {
	var res = ErrorResponseSimple{}
	if err != nil {
		res.Error = try.Try(err).Error()
	} else {
		res.Error = http.StatusText(http.StatusNotFound)
	}
	ctx.Response().Header().Add("error", res.Error)
	ctx.JSON(http.StatusNotFound, res)

	return errors.New(res.Error)
}

// OKResponse is 200 request
func (c Controller) OKResponse(ctx echo.Context, res interface{}) error {
	if res == nil {
		res = NormalResponse{}
	}
	ctx.JSON(http.StatusOK, res)

	return nil
}

// MustGetPayload is for payload middleware
func (c Controller) MustGetPayload(ctx echo.Context) interface{} {
	obj, ok := middlewares.GetPayload(ctx)
	assert.True(ok, "[BUG] payload un-marshaller failed")

	return obj
}
