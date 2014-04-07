package views

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dancannon/gonews/core/lib"
	"net/http"

	"github.com/dancannon/gonews/core/models"
)

type SecurityLogin struct {
	Errors ViewErrors `form:"-"`

	User       *models.User `form:"-"`
	Username   string       `form:"username"`
	Password   string       `form:"password"`
	RememberMe bool         `form:"remember-me"`
}

func (v *SecurityLogin) Validate(errors *binding.Errors, req *http.Request) {
	// Dont bother validating if there are already errors
	if errors.Count() > 0 {
		return
	}

	// Check username/password are correct
	user, err := lib.CheckUserLogin(v.Username, v.Password)
	if err != nil {
		errors.Overall["error"] = err.Error()
	}
	v.User = user
}
