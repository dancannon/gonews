package views

import (
	"github.com/dancannon/gonews/core/models"
)

type RulesNew struct {
	Id    string       `json:"id"`
	IsNew bool         `json:"is_new"`
	Post  string       `json:"post"`
	Rule  *models.Rule `json:"rule" binding:"required"`
}
