package infrastructure

import (
	"os"

	"github.com/dancannon/gonews/core/config"
	r "github.com/dancannon/gorethink"
)

var (
	rdb_session *r.Session
)

func InitRethinkDB(conf config.RethinkDB) {
	var err error

	address := conf.Address

	// Check for environment variables
	envHost := os.Getenv("RETHINKDB_HOST")
	envPort := os.Getenv("RETHINKDB_HOST")

	if envHost != "" && envPort != "" {
		address = envHost + ":" + envPort
	}

	rdb_session, err = r.Connect(map[string]interface{}{
		"address":  address,
		"database": conf.Database,
		"retries":  5,
	})
	if err != nil {
		panic(err)
	}
}

func RethinkDB() *r.Session {
	return rdb_session
}
