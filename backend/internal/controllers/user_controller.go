package controllers

import "backend/internal/services"

type UserController struct {
	userService *services.UserService
	jwtService  *services.JwtService
}

func (c *UserController) LogIn(username string, password string) (string, error) {
	err := c.userService.VerifyUser(username, password)
	if err != nil {
		return "", err
	}

	token, err := c.jwtService.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *UserController) LogInWithToken(token string) (string, error) {
	username, err := c.jwtService.ParseToken(token)
	if err != nil {
		return "", err
	}

	err = c.userService.DoesUserExist(username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (c *UserController) SignUp(username string, password string) (string, error) {
	err := c.userService.SignUp(username, password)
	if err != nil {
		return "", err
	}

	token, err := c.jwtService.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}
