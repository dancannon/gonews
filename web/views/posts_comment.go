package views

import (
	"github.com/dancannon/gonews/core/models"
	"github.com/russross/blackfriday"
	"html/template"
)

type PostsComment struct {
	Comment  models.Comment
	Children []PostsComment
}

func (v *PostsComment) HasChildren() bool {
	return len(v.Children) > 0
}

func (v *PostsComment) CommentContent() template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(v.Comment.Content))))
}
