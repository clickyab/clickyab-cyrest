package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"

	"github.com/Sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/labstack/echo.v3"
)

func (pl *changePasswordPayload) Validate(ctx echo.Context) error {
	errs := validator.New().Struct(pl)
	if errs == nil {
		return nil
	}
	res := middlewares.GroupError{}
	for _, i := range errs.(validator.ValidationErrors) {
		switch i.Field() {
		case "OldPassword":
			res["old_password"] = "old password is wrong"

		case "NewPassword":
			res["new_password"] = "new password can not be less than 6 charachter"

		default:
			logrus.Panicf("the field %s is not translated", i)
		}
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func (pl *loginPayload) Validate(ctx echo.Context) error {
	errs := validator.New().Struct(pl)
	if errs == nil {
		return nil
	}
	res := middlewares.GroupError{}
	for _, i := range errs.(validator.ValidationErrors) {
		switch i.Field() {
		case "Email":
			res["email"] = "email is invalid"

		case "Password":
			res["password"] = "password is too short"

		default:
			logrus.Panicf("the field %s is not translated", i)
		}
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func (pl *registrationPayload) Validate(ctx echo.Context) error {
	errs := validator.New().Struct(pl)
	if errs == nil {
		return nil
	}
	res := middlewares.GroupError{}
	for _, i := range errs.(validator.ValidationErrors) {
		switch i.Field() {
		case "Email":
			res["email"] = "invalid value"

		case "Password":
			res["password"] = "invalid value"

		default:
			logrus.Panicf("the field %s is not translated", i)
		}
	}

	if len(res) > 0 {
		return res
	}

	return nil
}
