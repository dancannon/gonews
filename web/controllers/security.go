package controllers

import (
	"github.com/dancannon/gonews/core/repos"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/gorilla/securecookie"

	"github.com/dancannon/gonews/core/lib"
	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/web/views"
)

// If user is not logged in then return an error page.
func SecuritySecure(w http.ResponseWriter, r *http.Request, user models.User) {
	if !user.Authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// If user is already logged in then redirect to homepage
func SecurityNotSecure(w http.ResponseWriter, r *http.Request, user models.User) {
	if user.Authenticated {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func SecurityLogin(
	c martini.Context,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	session sessions.Session,
	sc *securecookie.SecureCookie,
	viewData views.SecurityLogin,
	errors binding.Errors,
	user models.User,
) {
	if r.Method == "POST" {
		if errors.Count() == 0 {
			// Save user in session
			session.Set("user_id", viewData.User.Id)

			// Save user in martini context
			viewData.User.Authenticated = true
			c.Map(viewData.User)

			// If user choose "remember me" store user token in cookies
			if viewData.RememberMe {
				// Create cookie
				cookie, err := lib.CreateUserToken(user.Username, sc)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, cookie)
			}

			// If referrer was stored in session then redirect to that page otherwise
			// redirect to homepage
			if session.Get("security_target_path") != nil {
				referrer := session.Get("security_target_path").(string)
				session.Delete("security_target_path")

				http.Redirect(w, r, referrer, http.StatusTemporaryRedirect)
			} else {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			}

			return
		} else {
			// Add errors to view
			viewData.Errors = views.ViewErrors(errors)
		}
	}

	renderer.HTML(http.StatusOK, "security/login", viewData)
}

func SecurityLogout(c martini.Context, w http.ResponseWriter, r *http.Request, session sessions.Session, user models.User) {
	if user.Authenticated {
		// Delete session info if set
		if user_id := session.Get("user_id"); user_id != nil {
			session.Delete("user_id")
		}
		// Delete cookie if set
		if cookie, err := r.Cookie("security_token"); err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}

		c.Map(models.User{})
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func SecurityRegister(
	c martini.Context,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	viewData views.SecurityRegister,
	session sessions.Session,
	errors binding.Errors,
	user models.User,
) {
	if r.Method == "POST" {
		if errors.Count() == 0 {
			// Create new user document
			user, err := models.NewUser(viewData.Username, viewData.Email, viewData.Password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = repos.Users.Insert(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Login the user
			session.Set("user_id", user.Id)

			// If referrer was stored in session then redirect to that page otherwise
			// redirect to homepage
			if session.Get("security_target_path") != nil {
				referrer := session.Get("security_target_path").(string)
				session.Delete("security_target_path")

				http.Redirect(w, r, referrer, http.StatusTemporaryRedirect)
			} else {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			}
			return
		} else {
			// Add errors to view
			viewData.Errors = views.ViewErrors(errors)
		}
	}

	renderer.HTML(http.StatusOK, "security/register", viewData)
}
