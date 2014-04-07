package models

import (
	"time"

	"github.com/dancannon/gofetch/config"
)

type Rule struct {
	Id          string        `json:"id,omitempty", gorethink:"id,omitempty"`
	Name        string        `json:"name" binding:"required"`
	Type        string        `json:"type" binding:"required"`
	Host        string        `json:"host" binding:"required"`
	PathPattern string        `json:"path_pattern" binding:"required"`
	Url         string        `json:"test_url" binding:"required"`
	Values      []interface{} `json:"values"`
	AuthorId    string
	AuthorName  string
	Likes       int
	Dislikes    int
	Created     time.Time
	Modified    time.Time
}

func (r *Rule) ToGofetchRule() config.Rule {
	return config.Rule{
		Id:          r.Id,
		Type:        r.Type,
		Host:        r.Host,
		PathPattern: r.PathPattern,
		Values:      r.Values,
	}
}
