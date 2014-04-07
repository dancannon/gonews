package models

type Vote struct {
	Id     string `gorethink:"id,omitempty"`
	Entity string
	User   string
	Type   string
}
