package views

import (
	"github.com/dancannon/gonews/core/models"
	"github.com/russross/blackfriday"
	"html/template"
)

type CommentView struct {
	models.Comment
	Children []CommentView

	UserVote string
}

func (v *CommentView) RenderContent() template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(v.Comment.Content))))
}
