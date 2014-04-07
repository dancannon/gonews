package middleware

import (
	"github.com/codegangsta/martini"
	"github.com/dancannon/gonews/core/config"
)

func Config(conf *config.Config) martini.Handler {
	return func(c martini.Context) {
		c.Map(conf)
		c.Next()
	}
}

func ConfigFile(file string) martini.Handler {
	var conf *config.Config

	err := config.LoadFile(conf, file)
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		c.Map(conf)
		c.Next()
	}
}
