package views

import (
	"math"
)

type PostsList struct {
	Posts       []PostsView
	CurrentPage int
	PageCount   int
	TotalCount  int
}

func (l PostsList) TotalPages() int {
	return int(math.Ceil(float64(l.TotalCount) / float64(l.PageCount)))
}

func (l PostsList) PrevPage() int {
	if l.CurrentPage > 1 {
		return l.CurrentPage - 1
	}

	return 1
}

func (l PostsList) NextPage() int {
	return l.CurrentPage + 1
}
