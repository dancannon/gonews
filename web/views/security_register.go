package views

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dancannon/gonews/core/repos"
	"net/http"
)

type SecurityRegister struct {
	Errors ViewErrors `form:"-"`

	Username        string `form:"username" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"password_confirm" binding:"required"`
}

func (v SecurityRegister) Validate(errors *binding.Errors, req *http.Request) {
	// Dont bother validating if there are already errors
	if errors.Count() > 0 {
		return
	}

	// Check that a user does not exist with the same username
	uuCheck, err := repos.Users.FindByUsername(v.Username)
	if err != nil {
		errors.Overall["error"] = "An error occurred when creating your user account, please try again later."
		return
	}
	// Check that a user does not exist with the same email address
	ueCheck, err := repos.Users.FindByEmail(v.Email)
	if err != nil {
		errors.Overall["error"] = "An error occurred when creating your user account, please try again later."
		return
	}

	if uuCheck != nil {
		errors.Fields["Username"] = "A user with that username already exists"
	}
	if ueCheck != nil {
		errors.Fields["Email"] = "A user with that email already exists"
	}

	if v.Password != v.PasswordConfirm {
		errors.Fields["Password"] = "The passwords do not match"
	}
}
