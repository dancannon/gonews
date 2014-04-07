package views

import (
	"github.com/dancannon/gonews/core/models"
)

type RulesView struct {
	*models.Rule
	Comments []CommentView

	UserVote string
}
