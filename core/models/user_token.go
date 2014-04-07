package models

type UserToken struct {
	Id       string `gorethink:"id,omitempty"`
	Username string
	Secret   string
}
