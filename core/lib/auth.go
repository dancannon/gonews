package lib

import (
	"fmt"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/securecookie"
	"github.com/jmcvetta/randutil"

	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
)

func CheckUserLogin(username, password string) (user *models.User, err error) {
	user, err = repos.Users.FindByUsername(username)
	if err != nil || user == nil {
		err = fmt.Errorf("Incorrect username or password")
	} else if bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)) != nil {
		err = fmt.Errorf("Incorrect username or password")
	}

	return
}

func CreateUserToken(username string, sc *securecookie.SecureCookie) (*http.Cookie, error) {
	// Create cookie
	secret, err := randutil.AlphaString(64)
	if err != nil {
		return nil, err
	}

	// Create new token
	token := models.UserToken{
		Username: username,
		Secret:   secret,
	}

	// Save token in DB
	err = repos.UserTokens.Insert(&token)
	if err != nil {
		return nil, err
	}

	// Save cookie
	encoded, err := sc.Encode("security_token", token)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:  "security_token",
		Value: encoded,
		Path:  "/",
	}, nil
}
