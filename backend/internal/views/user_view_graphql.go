package views

import (
	"backend/internal/controllers"
	"backend/internal/middlewares"
	"context"
)

type UserViewGraphQl struct {
	controller         *controllers.UserController
	loginTokenName     string
	loginTokenMaxAge_s int
}

func (v *UserViewGraphQl) LogIn(
	ctx context.Context,
	username string,
	password string,
) (*bool, error) {
	var ret bool
	token, err := v.controller.LogIn(username, password)
	if err != nil {
		ret = false
		return &ret, err
	}

	setter := middlewares.GetCookieSetter(ctx)
	setter(v.loginTokenName, token, v.loginTokenMaxAge_s)
	ret = true
	return &ret, nil
}
