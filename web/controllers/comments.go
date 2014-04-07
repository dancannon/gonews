package controllers

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dancannon/gonews/core/config"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
	"github.com/dancannon/gonews/web/views"
)

func CommentsNew(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	errors binding.Errors,
	comment models.Comment,
	user models.User,
	conf *config.Config,
) {
	viewData := views.CommentsNew{
		Comment: comment,
	}

	if r.Method == "POST" {
		var post *models.Post
		var parent *models.Comment

		// Check that the post exists
		if comment.Post != "" {
			var err error
			post, err = repos.Posts.FindById(comment.Post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if post == nil {
				errors.Overall["post"] = "The parent post does not exist"
			}
		}

		// If the comment has a parent check it exists
		depth := 0
		if comment.Parent != "" {
			var err error
			parent, err = repos.Comments.FindById(comment.Parent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if parent == nil {
				errors.Overall["parent"] = "That parent comment does not exist"
			} else {
				depth = parent.Depth + 1
			}
		}

		if errors.Count() == 0 {
			// Add extra post info
			comment.AuthorId = user.Id
			comment.AuthorName = user.Username
			comment.Depth = depth
			comment.Created = time.Now()
			comment.Modified = time.Now()

			err := repos.Comments.Store(&comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/post/"+post.ActualUrl()+"#comment-"+comment.Id, http.StatusTemporaryRedirect)
			return
		} else {
			viewData.Errors = views.ViewErrors(errors)
		}
	}

	// Override the post/parent
	if postId := r.URL.Query().Get("post"); postId != "" {
		viewData.Comment.Post = postId
	}
	if parentId := r.URL.Query().Get("parent"); parentId != "" {
		viewData.Comment.Parent = parentId
	}

	// Ensure that the post ID is set
	if viewData.Post == "" {
		http.Error(w, "You must specify a post", http.StatusBadRequest)
		return
	}

	renderer.HTML(http.StatusOK, "comments/new", viewData)
}

func CommentsLikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	commentsVote(models.VoteTypeLike, params, w, r, renderer, user)
}

func CommentsDislikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	commentsVote(models.VoteTypeDislike, params, w, r, renderer, user)
}

func commentsVote(
	voteType string,
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	comment, err := repos.Comments.FindById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if comment == nil {
		http.NotFound(w, r)
		return
	}

	// Get a previous vote if one exists
	vote, err := repos.Votes.FindByEntityAndUser(comment.Id, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if vote != nil {
		if voteType != vote.Type {
			// If the vote was different change the value
			repos.Comments.UpdateVote(comment.Id, voteType, false, true)
			repos.Votes.UpdateFields(vote.Id, map[string]interface{}{
				"Type": voteType,
			})
		} else {
			// Otherwise delete the vote
			repos.Comments.UpdateVote(comment.Id, voteType, true, false)
			repos.Votes.Delete(vote.Id)
		}
	} else {
		repos.Comments.UpdateVote(comment.Id, voteType, false, false)
		repos.Votes.Store(&models.Vote{
			Entity: comment.Id,
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

func commentsLoad(user models.User, postId string) []views.CommentView {
	var comments []views.CommentView

	cs, err := repos.Comments.FindByPost(postId)
	if err != nil {
		return comments
	}

	for _, c := range cs {
		voteType := ""

		// Check if the user has voted for this comment
		vote, err := repos.Votes.FindByEntityAndUser(c.Id, user.Id)
		if err == nil && vote != nil {
			voteType = vote.Type
		}

		comments = append(comments, views.CommentView{
			Comment:  c,
			Children: commentsLoadChildren(user, c.Id),
			UserVote: voteType,
		})
	}

	return comments
}

func commentsLoadChildren(user models.User, commentId string) []views.CommentView {
	var comments []views.CommentView

	cs, err := repos.Comments.FindChildren(commentId)
	if err != nil {
		return comments
	}

	for _, c := range cs {
		voteType := ""

		// Check if the user has voted for this comment
		vote, err := repos.Votes.FindByEntityAndUser(c.Id, user.Id)
		if err == nil && vote != nil {
			voteType = vote.Type
		}

		comments = append(comments, views.CommentView{
			Comment:  c,
			Children: commentsLoadChildren(user, c.Id),
			UserVote: voteType,
		})
	}

	return comments
}
