package models

import (
	"github.com/codegangsta/martini-contrib/binding"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	PostTypeLink = "link"
	PostTypeText = "text"

	LinkTypeArticle = "article"

	VoteTypeLike    = "like"
	VoteTypeDislike = "dislike"
)

type Post struct {
	Id           string `gorethink:"id,omitempty"`
	Url          string
	SlugId       string
	Slug         string
	AuthorId     string
	AuthorName   string
	Type         string `form:"type" binding:"required"`
	Title        string `form:"title" binding:"required"`
	Link         string `form:"link"`
	Content      string `form:"content"`
	EmbedRule    string
	EmbedType    string
	EmbedContent map[string]interface{}
	Likes        int
	Dislikes     int
	Created      time.Time
	Modified     time.Time
}

func (p *Post) Validate(errors *binding.Errors, req *http.Request) {
	if p.Type == PostTypeLink {
		if p.Link == "" {
			errors.Fields["Link"] = "Required"
		}
		if u, err := url.Parse(p.Link); err != nil {
			errors.Fields["Link"] = "The link must be a valid URL"
		} else {
			// Ensure URL is valid for extraction
			if u.Scheme == "" {
				u.Scheme = "http"
				p.Link = u.String()
			}
		}
	} else if p.Type == PostTypeText {
		if p.Content == "" {
			errors.Fields["Content"] = "Required"
		}
	}
}

func (p *Post) Score() string {
	return strconv.Itoa(p.Likes - p.Dislikes)
}

func (p *Post) Rank() float64 {
	var score = float64(p.Likes - p.Dislikes)
	var order = math.Log10(math.Max(math.Abs(score), 1))
	var sign int64
	if score < 1 {
		sign = -1
	} else {
		sign = 1
	}
	var seconds = p.Created.Unix() - 1134028003

	return (order + float64((sign*seconds)/45000))
}

func (p *Post) ActualUrl() string {
	if p.Url != "" {
		return p.Url
	} else {
		return "view/" + p.Id
	}
}

func (p *Post) Host() string {
	if p.Type == "link" {
		u, err := url.Parse(p.Link)
		if err != nil {
			return ""
		}
		return u.Host
	}

	return ""
}

func (p *Post) Path() string {
	if p.Type == "link" {
		u, err := url.Parse(p.Link)
		if err != nil {
			return ""
		}
		return u.RequestURI()
	}

	return ""
}
