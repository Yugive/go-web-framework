package main

import (
	"go-web-framework/framework"
)

func UserLoginController(c *framework.Context) error {
	c.Json("ok, userController 444")
	return nil
}
