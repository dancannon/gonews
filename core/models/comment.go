package models

import (
	"time"
)

type Comment struct {
	Id         string `gorethink:"id,omitempty"`
	Post       string `form:"post" binding:"required"`
	Parent     string `form:"parent"`
	Depth      int
	AuthorId   string
	AuthorName string
	Content    string `form:"content" binding:"required"`
	Likes      int
	Dislikes   int
	Created    time.Time
	Modified   time.Time
}
