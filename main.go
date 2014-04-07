package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dancannon/gonews/core/config"
	"github.com/dancannon/gonews/core/infrastructure"
	"github.com/dancannon/gonews/data"
	"github.com/dancannon/gonews/web"
)

const (
	VERSION = "0.5.0"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	configPath := flag.String("config", "config.toml", "Config file.")
	env := flag.String("env", "", "Enviroment, overrides config.")
	version := flag.Bool("version", false, "Output version and exit")
	setupDb := flag.Bool("setup-db", false, "Setup database")
	exampleData := flag.Bool("example-data", false, "Insert example data")

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Load config
	conf := config.NewConfig()
	err := config.LoadFile(conf, *configPath)
	if err != nil {
		log.Fatal("Error reading config: ", err)
	}

	// Override env if needed
	if *env != "" {
		conf.Env = *env
	}

	// Various package initialization
	rand.Seed(time.Now().UnixNano())
	infrastructure.InitRedis(conf.Redis)
	infrastructure.InitRethinkDB(conf.RethinkDB)

	// Setup DB if specified
	if *setupDb {
		log.Println("Setting up database")
		data.Setup(*conf, *exampleData)
	}

	// Create and start the server
	go func() {
		server := web.NewServer(conf)
		server.Run()
	}()

	// Wait for exit
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
