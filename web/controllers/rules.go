package controllers

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/dancannon/gofetch"
	gfc "github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	"github.com/dancannon/gonews/core/config"
	"github.com/dancannon/gonews/core/lib"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"
	"github.com/dancannon/gonews/web/views"
)

func RulesList(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	// Load rules
	var err error
	var rules []*models.Rule

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

	post := r.URL.Query().Get("post")
	query := r.URL.Query().Get("q")
	host, ok := params["host"]

	if ok && host != "" {
		if query != "" {
			rules, err = repos.Rules.FindByQueryHostAndPage(r.URL.Query().Get("q"), host, page, count)
		} else {
			rules, err = repos.Rules.FindByHostAndPage(host, page, count)
		}
	} else {
		if query != "" {
			rules, err = repos.Rules.FindByQueryAndPage(r.URL.Query().Get("q"), page, count)
		} else {
			rules, err = repos.Rules.FindByPage(page, count)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewData := views.RulesList{
		Query:       query,
		Post:        post,
		Host:        host,
		CurrentPage: page,
		PageCount:   count,
		TotalCount:  len(rules),
	}

	// Check user votes for each rule
	ruleViews := make([]views.RulesView, len(rules))
	for i, rule := range rules {
		ruleView := views.RulesView{
			Rule: rule,
		}

		vote, err := repos.Votes.FindByEntityAndUser(rule.Id, user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if vote != nil {
			ruleView.UserVote = vote.Type
		}

		ruleViews[i] = ruleView
	}
	viewData.Rules = ruleViews

	renderer.HTML(http.StatusOK, "rules/list", viewData)
}

func RulesNew(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
) {
	var viewData views.RulesNew
	var post *models.Post

	postId := r.URL.Query().Get("post")

	if postId != "" {
		post, err := repos.Posts.FindById(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if post == nil {
			http.NotFound(w, r)
			return
		}
	}

	// Load the rule if the ID was specified
	if id, ok := params["id"]; ok {
		rule, err := repos.Rules.FindById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rule == nil {
			http.NotFound(w, r)
			return
		}

		viewData = views.RulesNew{
			Id:    id,
			Rule:  rule,
			Post:  postId,
			IsNew: false,
		}
	} else {
		if post != nil {
			viewData = views.RulesNew{
				IsNew: true,
				Post:  postId,
				Rule: &models.Rule{
					Host: post.Host(),
					Url:  post.Url,
				},
			}
		} else {
			viewData = views.RulesNew{
				IsNew: true,
				Post:  postId,
			}
		}
	}

	renderer.HTML(http.StatusOK, "rules/new", viewData)
}

func RulesSave(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	viewData views.RulesNew,
	conf *config.Config,
	user models.User,
) {
	// If ID was specified then load from the database
	var err error
	var rule *models.Rule

	if viewData.Rule.Id != "" {
		rule, err = repos.Rules.FindById(viewData.Rule.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rule == nil {
			http.NotFound(w, r)
			return
		}
	}

	// Validate rule
	c, err := gfc.LoadConfig(conf.GoFetch.ConfigFile)
	if err != nil {
		renderer.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.Rules = []gfc.Rule{viewData.Rule.ToGofetchRule()}

	fetcher := gofetch.NewFetcher(c)
	res, err := fetcher.Fetch(viewData.Rule.Url)
	if err != nil {
		renderer.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Add extra rule info
	viewData.Rule.AuthorId = user.Id
	viewData.Rule.AuthorName = user.Username
	viewData.Rule.Created = time.Now()
	viewData.Rule.Modified = time.Now()

	if viewData.Rule.Id != "" {
		err = repos.Rules.Update(viewData.Rule)
	} else {
		err = repos.Rules.Store(viewData.Rule)
	}
	if err != nil {
		renderer.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Refresh scraped post content if post ID was set
	if viewData.Post != "" {
		go func() {
			post, err := repos.Posts.FindById(viewData.Post)
			if err != nil {
				return
			}

			lib.ExtractLinkContentWithRule(conf.GoFetch, *post, *viewData.Rule)
		}()
	}

	renderer.JSON(http.StatusOK, map[string]interface{}{
		"id":       viewData.Rule.Id,
		"response": res,
	})
}

func RulesTest(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	viewData views.RulesNew,
	conf *config.Config,
) {
	// Validate rule
	c, err := gfc.LoadConfig(conf.GoFetch.ConfigFile)
	if err != nil {
		renderer.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.Rules = []gfc.Rule{viewData.Rule.ToGofetchRule()}

	fetcher := gofetch.NewFetcher(c)
	res, err := fetcher.Fetch(viewData.Rule.Url)
	if err != nil {
		renderer.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	renderer.JSON(http.StatusOK, map[string]interface{}{
		"response": res,
	})
}

func RulesLikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	rulesVote(models.VoteTypeLike, params, w, r, renderer, user)
}

func RulesDislikeVote(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	rulesVote(models.VoteTypeDislike, params, w, r, renderer, user)
}

func rulesVote(
	voteType string,
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
	user models.User,
) {
	rule, err := repos.Rules.FindById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rule == nil {
		http.NotFound(w, r)
		return
	}

	// Get a previous vote if one exists
	vote, err := repos.Votes.FindByEntityAndUser(rule.Id, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if vote != nil {
		if voteType != vote.Type {
			// If the vote was different change the value
			repos.Rules.UpdateVote(rule.Id, voteType, false, true)
			repos.Votes.UpdateFields(vote.Id, map[string]interface{}{
				"Type": voteType,
			})
		} else {
			// Otherwise delete the vote
			repos.Rules.UpdateVote(rule.Id, voteType, true, false)
			repos.Votes.Delete(vote.Id)
		}
	} else {
		repos.Rules.UpdateVote(rule.Id, voteType, false, false)
		repos.Votes.Store(&models.Vote{
			Entity: rule.Id,
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

func RulesLoadUrl(
	params martini.Params,
	w http.ResponseWriter,
	r *http.Request,
	renderer render.Render,
) {
	// Get URL
	if r.URL.Query().Get("url") == "" {
		http.NotFound(w, r)
		return
	}

	u, err := url.Parse(r.URL.Query().Get("url"))
	if err != nil {
		http.Error(w, "The link must be a valid URL", http.StatusBadRequest)
	} else {
		// Ensure URL is valid for extraction
		if u.Scheme == "" {
			u.Scheme = "http"
		}
	}

	// Fetch page
	response, err := http.Get(u.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Ensure that the document is a HTML page
	if !strings.Contains(response.Header.Get("Content-Type"), "html") {
		http.Error(w, "Invalid page type", http.StatusBadRequest)
		return
	}

	// Create document
	doc, err := document.NewDocument(u.String(), response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert any relative URLs to absolute URLs
	convertUrls(u.String(), doc.Doc.Node())
	removeScripts(doc.Doc.Node())

	// Add JS to document
	ns, _ := html.ParseFragment(strings.NewReader(`
		<style>
			.highlight{
				background-color: #bcd5eb !important;
			}
			.highlight-selected{
				background-color: #bcd5eb !important;
				outline: 2px solid #5166bb !important;
			}
		</style>
		<script>
			(function () {
				var prev;
				var closestAnchor = function(el) {
					if (el.tagName == "A") {
						return el;
					} else if (el.parentElement != null) {
						return closestAnchor(el.parentElement);
					} else {
						return null;
					}
				};
				document.addEventListener('mouseout', function (e) {
					event.target.classList.remove('highlight');
					prev = undefined;
				});
				document.addEventListener('mouseover', function (e) {
					if (event.target === document.body ||
						(prev && prev === event.target)) {
						return;
					}
					if (prev) {
						prev.classList.remove('highlight');
						prev = undefined;
					}
					if (event.target) {
						prev = event.target;
						prev.classList.add('highlight');
					}
				});
				document.addEventListener('click', function (e) {
					if (!e.ctrlKey) {
						e.target.classList.toggle('highlight-selected');
						e.preventDefault();
					} else {
						var link = closestAnchor(e.target);
						if (link != null && link.getAttribute("href") != null) {
							window.location = "/rule/load_url?url="+link.getAttribute("href");
						}

						e.preventDefault();
					}

					return false;
				});
				window.parent.updateTestUrl("`+u.String()+`");
			})();
		</script>
	`), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	for _, n := range ns {
		n.Parent = nil
		n.PrevSibling = nil
		n.NextSibling = nil

		doc.Body.Node().AppendChild(&(*n))
	}

	// Write document to response writer
	html.Render(w, doc.Doc.Node())
}

func convertUrls(u string, n *html.Node) {
	if n.Type == html.ElementNode {
		// Ensure that the body tag is added to the result document
		tmpAttrs := []html.Attribute{}
		for _, a := range n.Attr {
			if a.Key == "href" || a.Key == "src" {
				// Attempt to fix URLs
				urlb, err := url.Parse(u)
				if err != nil {
					continue
				}
				ur, err := url.Parse(a.Val)
				if err != nil {
					continue
				}
				a.Val = urlb.ResolveReference(ur).String()
			}

			tmpAttrs = append(tmpAttrs, a)
		}
		n.Attr = tmpAttrs
	}

	// Build the list of children node before iterating. This is needed because we will be
	// deleting nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		convertUrls(u, c)
	}
}

func removeScripts(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "script" {
		deleteNode(n)
	}

	// Build the list of children node before iterating. This is needed because we will be
	// deleting nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeScripts(c)
	}
}

func deleteNode(n *html.Node) {
	if n.Parent != nil {

	}
	if n.Parent.FirstChild == n {
		n.FirstChild = n.NextSibling
	}
	if n.NextSibling != nil {
		n.NextSibling.PrevSibling = n.PrevSibling
	}
	if n.Parent.LastChild == n {
		n.LastChild = n.PrevSibling
	}
	if n.PrevSibling != nil {
		n.PrevSibling.NextSibling = n.NextSibling
	}
	n.Parent = nil
}
