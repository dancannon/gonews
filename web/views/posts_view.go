package views

import (
	"github.com/dancannon/gonews/core/models"
	"github.com/russross/blackfriday"
	"html/template"
)

type PostsView struct {
	*models.Post
	Comments []CommentView

	UserVote string
}

func (v PostsView) RenderContent() template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(v.Post.Content))))
}

func (p PostsView) CreateCommentsNewView() CommentsNew {
	return CommentsNew{
		Comment: models.Comment{
			Post: p.Id,
		},
	}
}
