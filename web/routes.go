package web

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dancannon/gonews/core/models"

	"github.com/dancannon/gonews/web/controllers"
	"github.com/dancannon/gonews/web/views"
)

func initRouter() martini.Router {
	router := martini.NewRouter()

	// General pages
	// router.Any("/", controllers.GeneralHomepage)
	router.Any("/secret", controllers.SecuritySecure, controllers.GeneralSecret)

	// Security
	router.Any("/login", controllers.SecurityNotSecure, binding.Form(views.SecurityLogin{}), controllers.SecurityLogin)
	router.Any("/logout", controllers.SecurityLogout)
	router.Any("/register", controllers.SecurityNotSecure, binding.Form(views.SecurityRegister{}), controllers.SecurityRegister)

	// Posts
	router.Any("/", controllers.PostsPopularList)
	router.Any("/popular", controllers.PostsPopularList)
	router.Any("/top", controllers.PostsTopList)
	router.Any("/new", controllers.PostsNewList)
	router.Any("/posts/new", controllers.SecuritySecure, binding.Form(models.Post{}), controllers.PostsNew)
	router.Any("/post/view/:id", controllers.PostsIdView)
	router.Any("/post/like/:id", controllers.SecuritySecure, controllers.PostsLikeVote)
	router.Any("/post/refresh/:id", controllers.SecuritySecure, controllers.PostsRefresh)
	router.Any("/post/refresh/:id/:rule", controllers.SecuritySecure, controllers.PostsRefresh)
	router.Any("/post/dislike/:id", controllers.SecuritySecure, controllers.PostsDislikeVote)
	router.Any("/post/**", controllers.PostsUrlView)

	// Comments
	router.Any("/comments/new", controllers.SecuritySecure, binding.Form(models.Comment{}), controllers.CommentsNew)
	router.Any("/comment/like/:id", controllers.SecuritySecure, controllers.CommentsLikeVote)
	router.Any("/comment/dislike/:id", controllers.SecuritySecure, controllers.CommentsDislikeVote)

	// Rule Builder
	router.Any("/rules", controllers.RulesList)
	router.Any("/rules/host/:host", controllers.RulesList)
	router.Any("/rules/new", controllers.SecuritySecure, controllers.RulesNew)
	router.Any("/rule/save", controllers.SecuritySecure, binding.Json(views.RulesNew{}), controllers.RulesSave)
	router.Any("/rule/test", binding.Json(views.RulesNew{}), controllers.RulesTest)
	router.Any("/rule/edit/:id", controllers.SecuritySecure, controllers.RulesNew)
	router.Any("/rule/like/:id", controllers.SecuritySecure, controllers.RulesLikeVote)
	router.Any("/rule/dislike/:id", controllers.SecuritySecure, controllers.RulesDislikeVote)
	router.Any("/rule/load_url", controllers.SecuritySecure, controllers.RulesLoadUrl)

	return router
}
