package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/labstack/echo.v3"
)

func (pl registrationPayload) Validate(ctx echo.Context) error {
	return validator.New().Struct(pl)
}
