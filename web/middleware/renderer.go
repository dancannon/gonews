package middleware

import (
	"html/template"
	"strings"

	"github.com/dancannon/gonews/core/config"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/dancannon/gonews/core/models"
)

// Add more helper functions to the renderer
func RendererFuncs() martini.Handler {
	return func(
		c martini.Context,
		renderer render.Render,
		user models.User,
		conf *config.Config,
	) {
		renderer.Template().Funcs(template.FuncMap{
			"raw": func(input string) template.HTML {
				return template.HTML(input)
			},
			"nl2br": func(input string) template.HTML {
				return template.HTML(strings.Replace(input, "\n", "<br />", -1))
			},
			"is_authenticated": func() bool {
				return user.Authenticated
			},
			"get_user": func() models.User {
				return user
			},
			"get_env": func() string {
				return conf.Env
			},
		})

		c.Next()
	}
}

func BaseRendererFuncs() template.FuncMap {
	return template.FuncMap{
		"raw": func(input string) template.HTML {
			return template.HTML("")
		},
		"nl2br": func(input string) string {
			return ""
		},
		"is_authenticated": func() bool {
			return false
		},
		"get_user": func() models.User {
			return models.User{}
		},
		"get_env": func() string {
			return ""
		},
	}
}
