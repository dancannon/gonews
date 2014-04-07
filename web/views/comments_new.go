package views

import "github.com/dancannon/gonews/core/models"

type CommentsNew struct {
	Errors ViewErrors
	models.Comment
}
