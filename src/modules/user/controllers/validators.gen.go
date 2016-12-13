package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"

	"github.com/Sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/labstack/echo.v3"
)

func (pl changePasswordPayload) Validate(ctx echo.Context) error {
	errs := validator.New().Struct(pl)
	if errs == nil {
		return nil
	}
	res := middlewares.GroupError{}
	for i := range errs.(validator.ValidationErrors) {
		switch i {
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

func (pl registrationPayload) Validate(ctx echo.Context) error {
	errs := validator.New().Struct(pl)
	if errs == nil {
		return nil
	}
	res := middlewares.GroupError{}
	for i := range errs.(validator.ValidationErrors) {
		switch i {
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
