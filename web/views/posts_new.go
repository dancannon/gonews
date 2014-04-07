package views

import "github.com/dancannon/gonews/core/models"

type PostsNew struct {
	Errors ViewErrors
	models.Post
}
