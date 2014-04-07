package middleware

import (
	"github.com/codegangsta/martini-contrib/sessions"
	"net/http"

	"github.com/codegangsta/martini"
)

func AuthRedierct() martini.Handler {
	return func(
		c martini.Context,
		w http.ResponseWriter,
		r *http.Request,
		session sessions.Session,
	) {
		arw := authRedirectResponseWriter{r, session, w.(martini.ResponseWriter)}
		c.MapTo(arw, (*http.ResponseWriter)(nil))

		c.Next()
	}
}

type authRedirectResponseWriter struct {
	r       *http.Request
	session sessions.Session
	martini.ResponseWriter
}

func (w authRedirectResponseWriter) WriteHeader(code int) {
	if code == http.StatusUnauthorized || code == http.StatusForbidden {
		w.session.Set("security_target_path", w.r.URL.String())
		http.Redirect(w, w.r, "/login", http.StatusTemporaryRedirect)
		return
	}

	w.ResponseWriter.WriteHeader(code)
}
