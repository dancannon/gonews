package views

import (
	"math"
)

type RulesList struct {
	Rules       []RulesView
	Post        string
	Query       string
	Host        string
	CurrentPage int
	PageCount   int
	TotalCount  int
}

func (l RulesList) TotalPages() int {
	return int(math.Ceil(float64(l.TotalCount) / float64(l.PageCount)))
}

func (l RulesList) PrevPage() int {
	if l.CurrentPage > 1 {
		return l.CurrentPage - 1
	}

	return 1
}

func (l RulesList) NextPage() int {
	return l.CurrentPage + 1
}
