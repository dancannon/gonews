package middleware

import (
	"github.com/codegangsta/martini"
	"github.com/gorilla/securecookie"
)

func SecureCookie(hashKey, blockKey []byte) martini.Handler {
	var s = securecookie.New(hashKey, blockKey)

	return func(c martini.Context) {
		c.Map(s)
		c.Next()
	}
}
