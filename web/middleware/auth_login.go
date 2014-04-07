package middleware

import (
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/gorilla/securecookie"

	"github.com/dancannon/gonews/core/lib"
	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
)

func AuthLogin() martini.Handler {
	return func(
		c martini.Context,
		w http.ResponseWriter,
		r *http.Request,
		session sessions.Session,
		sc *securecookie.SecureCookie,
	) {
		var user models.User
		c.Map(user)

		// First check session
		if user_id := session.Get("user_id"); user_id != nil {
			var err error

			u, err := repos.Users.FindById(user_id.(string))
			if err != nil {
				session.Delete("user_id")

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if u == nil {
				session.Delete("user_id")
				return
			}

			user = *u
			user.Authenticated = true
			c.Map(user)
		} else if cookie, err := r.Cookie("security_token"); err == nil {
			// Otherwise check cookie
			var token models.UserToken

			if err = sc.Decode("security_token", cookie.Value, &token); err != nil {
				cookie.MaxAge = 0
				http.SetCookie(w, cookie)
				return
			}

			// Find token in the DB
			token, err := repos.UserTokens.FindToken(token.Username, token.Secret)
			if err != nil {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
				return
			}

			// Regenerate token
			err = repos.UserTokens.Delete(token.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cookie, err = lib.CreateUserToken(token.Username, sc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Fetch user
			u, err := repos.Users.FindByUsername(token.Username)
			if err != nil {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
				session.Delete("user_id")

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if u == nil {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
				session.Delete("user_id")
				return
			}

			session.Set("user_id", user.Id)
			http.SetCookie(w, cookie)
			user.Authenticated = true

			c.Map(user)
		}
	}
}
