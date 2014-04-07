package lib

import (
	"log"
	"net/url"
	"strings"

	"github.com/dancannon/gofetch"
	gfc "github.com/dancannon/gofetch/config"
	"github.com/dancannon/gonews/core/config"

	"github.com/extemporalgenome/slug"

	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
)

const (
	alphanum  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	minLength = 5
)

func GeneratePostUrl(post models.Post) {
	post.Slug = slug.Slug(post.Title)

	id := strings.Replace(post.Id, "-", "", -1)
	for i := minLength; i < len(id); i++ {
		post.SlugId = id[0:i]
		post.Url = post.SlugId + "/" + post.Slug

		// Check if a post exists with the same url
		post, err := repos.Posts.FindByUrl(post.Url)
		if err != nil {
			return
		}
		if post == nil {
			break
		}
	}

	repos.Posts.UpdateFields(post.Id, map[string]interface{}{
		"Slug":   post.Slug,
		"SlugId": post.SlugId,
		"Url":    post.Url,
	})
}

func ExtractLinkContent(conf config.GoFetch, post models.Post) {
	if post.Type == "link" {
		// Load config
		c, err := gfc.LoadConfig(conf.ConfigFile)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Parse the URL
		u, err := url.Parse(post.Link)
		if err == nil {
			// Load any relevant rules from the database
			rules, err := repos.Rules.FindByHost(u.Host)
			if err != nil {
				log.Println(err.Error())
				return
			}

			i := 1
			for _, rule := range rules {
				gfr := rule.ToGofetchRule()
				gfr.Priority = i
				c.Rules = append(c.Rules, gfr)
			}
		}

		fetcher := gofetch.NewFetcher(c)
		res, err := fetcher.Fetch(post.Link)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = repos.Posts.UpdateFields(post.Id, map[string]interface{}{
			"EmbedType":    res.PageType,
			"EmbedContent": res.Content,
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	return
}

func ExtractLinkContentWithRule(conf config.GoFetch, post models.Post, rule models.Rule) {
	if post.Type == "link" {
		// Load config
		c, err := gfc.LoadConfig(conf.ConfigFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.Rules = gfc.RuleSlice{rule.ToGofetchRule()}

		fetcher := gofetch.NewFetcher(c)
		res, err := fetcher.Fetch(post.Link)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = repos.Posts.UpdateFields(post.Id, map[string]interface{}{
			"EmbedRule":    rule.Id,
			"EmbedType":    res.PageType,
			"EmbedContent": res.Content,
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	return
}
