package data

import (
	"github.com/dancannon/gonews/core/config"
)

func Setup(conf config.Config, exampleData bool) {
	setupRethinkDB(conf, exampleData)
}
