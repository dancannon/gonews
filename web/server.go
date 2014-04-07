package web

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"

	"github.com/dancannon/gonews/core/config"
	"github.com/dancannon/gonews/web/middleware"
)

type Server struct {
	conf   *config.Config
	m      *martini.Martini
	Router martini.Router
}

func NewServer(conf *config.Config) *Server {
	martini.Env = conf.Env

	m := martini.New()

	// Load routes
	router := initRouter()

	// Add middleware
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static(conf.PublicPath))
	m.Use(render.Renderer(render.Options{
		Directory: conf.TemplatePath,
		Layout:    "base",
		Funcs:     []template.FuncMap{middleware.BaseRendererFuncs()},
	}))
	m.Use(middleware.Config(conf))
	m.Use(sessions.Sessions("session_id", sessions.NewCookieStore([]byte(conf.Security.Secret))))
	m.Use(middleware.SecureCookie([]byte(conf.Security.Secret), []byte(conf.Security.BlockKey)))
	m.Use(middleware.AuthLogin())
	m.Use(middleware.AuthRedierct())
	m.Use(middleware.RendererFuncs())

	m.Action(router.Handle)

	return &Server{conf, m, router}
}

func (s *Server) Run() {
	address := s.conf.Address
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host != "" || port != "" {
		address = host + ":" + port
	}

	log.Println("Starting server, listening on " + address)
	log.Fatalln(http.ListenAndServe(address, s.m))
}
