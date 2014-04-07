package models

import (
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

type Role string

const (
	ROLE_ADMIN string = "ADMIN"
)

type User struct {
	Id       string `gorethink:"id,omitempty"`
	Username string
	Email    string

	// Security values
	Hash string

	Active bool
	Roles  []string

	LastVisit time.Time
	Created   time.Time
	Modified  time.Time

	Authenticated bool `gorethink:"-"`
}

func NewUser(username, email, password string) (*User, error) {
	// Create new user document
	user := &User{}
	user.Username = username
	user.Email = email

	hp, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return nil, err
	}
	user.Hash = string(hp)

	return user, nil
}
