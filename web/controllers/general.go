package controllers

import (
	"net/http"

	"github.com/codegangsta/martini-contrib/render"
)

func GeneralHomepage(renderer render.Render) {
	renderer.HTML(http.StatusOK, "index", nil)
}

func GeneralSecret() string {
	return "Secret Page!"
}
