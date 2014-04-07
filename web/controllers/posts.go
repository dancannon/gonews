package controllers

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dancannon/gonews/core/config"
	"github.com/dancannon/gonews/core/lib"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
	"github.com/dancannon/gonews/web/views"
)

func PostsPopularList(
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	postsList("popular", w, r, renderer, user)
}

func PostsTopList(
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	postsList("top", w, r, renderer, user)
}

func PostsNewList(
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	postsList("new", w, r, renderer, user)
}

func postsList(
	sortMethod string,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	var err error
	var posts []*models.Post

	// Load pagination fields
	var page int
	var count int

	if r.URL.Query().Get("page") != "" {
		page, _ = strconv.Atoi(r.URL.Query().Get("page"))

		if page < 1 {
			page = 1
		}
	} else {
		page = 1
	}

	if r.URL.Query().Get("count") != "" {
		count, _ = strconv.Atoi(r.URL.Query().Get("count"))

		// Ensure that count is capped at 100
		if count > 100 {
			count = 100
		}
	} else {
		count = 25
	}

	if sortMethod == "top" {
		posts, err = repos.Posts.FindTopByPage(1, 25)
	} else if sortMethod == "new" {
		posts, err = repos.Posts.FindNewByPage(1, 25)
	} else {
		posts, err = repos.Posts.FindPopularByPage(1, 25)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewData := views.PostsList{
		CurrentPage: page,
		PageCount:   count,
		TotalCount:  len(posts),
	}

	// Check user votes for each post
	postsView := make([]views.PostsView, len(posts))
	for i, post := range posts {
		postView := views.PostsView{
			Post: post,
		}

		vote, err := repos.Votes.FindByEntityAndUser(post.Id, user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if vote != nil {
			postView.UserVote = vote.Type
		}

		postsView[i] = postView
	}
	viewData.Posts = postsView

	renderer.HTML(http.StatusOK, "posts/list", viewData)
}

func PostsIdView(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	post, err := repos.Posts.FindById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	postsView(post, params, w, r, renderer, user)
}

func PostsUrlView(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	post, err := repos.Posts.FindByUrl(params["_1"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	postsView(post, params, w, r, renderer, user)
}

func postsView(
	post *models.Post,
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	viewData := views.PostsView{
		Post:     post,
		Comments: commentsLoad(user, post.Id),
	}

	// Check if the user has voted for this post
	vote, err := repos.Votes.FindByEntityAndUser(post.Id, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if vote != nil {
		viewData.UserVote = vote.Type
	}

	renderer.HTML(http.StatusOK, "posts/view", viewData)
}

func PostsRefresh(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	conf *config.Config,
) {
	post, err := repos.Posts.FindById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	ruleId, ok := params["rule"]
	if ok {
		post.EmbedRule = ruleId
	}

	log.Println(post.EmbedRule)
	if post.EmbedRule != "" {
		// Fetch rule
		rule, err := repos.Rules.FindById(post.EmbedRule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go lib.ExtractLinkContentWithRule(conf.GoFetch, *post, *rule)
	} else {
		go lib.ExtractLinkContent(conf.GoFetch, *post)
	}

	http.Redirect(w, r, "/post/view/"+post.Id, http.StatusTemporaryRedirect)
}

func PostsNew(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	errors binding.Errors,
	post models.Post,
	user models.User,
	conf *config.Config,
) {
	viewData := views.PostsNew{
		Post: post,
	}

	if r.Method == "POST" {
		if errors.Count() == 0 {
			// Add extra post info
			post.AuthorId = user.Id
			post.AuthorName = user.Username
			post.EmbedType = "fetching"
			post.Created = time.Now()
			post.Modified = time.Now()

			err := repos.Posts.Store(&post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Start processing
			go lib.GeneratePostUrl(post)
			go lib.ExtractLinkContent(conf.GoFetch, post)

			http.Redirect(w, r, "/post/view/"+post.Id, http.StatusTemporaryRedirect)
			return
		} else {
			viewData.Errors = views.ViewErrors(errors)
		}
	}

	renderer.HTML(http.StatusOK, "posts/new", viewData)
}

func PostsLikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	postsVote(models.VoteTypeLike, params, w, r, renderer, user)
}

func PostsDislikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	postsVote(models.VoteTypeDislike, params, w, r, renderer, user)
}

func postsVote(
	voteType string,
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	post, err := repos.Posts.FindById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	// Get a previous vote if one exists
	vote, err := repos.Votes.FindByEntityAndUser(post.Id, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if vote != nil {
		if voteType != vote.Type {
			// If the vote was different change the value
			repos.Posts.UpdateVote(post.Id, voteType, false, true)
			repos.Votes.UpdateFields(vote.Id, map[string]interface{}{
				"Type": voteType,
			})
		} else {
			// Otherwise delete the vote
			repos.Posts.UpdateVote(post.Id, voteType, true, false)
			repos.Votes.Delete(vote.Id)
		}
	} else {
		repos.Posts.UpdateVote(post.Id, voteType, false, false)
		repos.Votes.Store(&models.Vote{
			Entity: post.Id,
			User:   user.Id,
			Type:   voteType,
		})
	}

	if r.Referer() == "" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, r.Referer(), http.StatusTemporaryRedirect)
	}
}
