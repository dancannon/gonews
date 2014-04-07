package data

import (
	"log"
	"time"

	"github.com/dancannon/gonews/core/config"
	"github.com/dancannon/gonews/core/lib"

	"github.com/dancannon/gonews/core/infrastructure"
	"github.com/dancannon/gonews/core/models"
	"github.com/dancannon/gonews/core/repos"

	r "github.com/dancannon/gorethink"
)

var (
	userAdmin string
	user1     string
	user2     string
)

func setupRethinkDB(conf config.Config, exampleData bool) {
	createDatabase()
	createTables()
	createIndexes()

	if exampleData {
		createUsers()
		createPosts(conf)
		createRules()
	}
}

func createDatabase() {
	r.DbDrop("news").RunWrite(infrastructure.RethinkDB())
	r.DbCreate("news").RunWrite(infrastructure.RethinkDB())
}

func createTables() {
	r.Db("news").TableDrop("users").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableDrop("user_tokens").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableDrop("rules").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableDrop("posts").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableDrop("comments").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableDrop("votes").RunWrite(infrastructure.RethinkDB())

	r.Db("news").TableCreate("users").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableCreate("user_tokens").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableCreate("rules").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableCreate("posts").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableCreate("comments").RunWrite(infrastructure.RethinkDB())
	r.Db("news").TableCreate("votes").RunWrite(infrastructure.RethinkDB())
}

func createIndexes() {
	r.Db("news").Table("users").IndexCreate("Username").RunWrite(infrastructure.RethinkDB())
	r.Db("news").Table("users").IndexCreate("Email").RunWrite(infrastructure.RethinkDB())
}

func createUsers() {
	var user *models.User
	var err error

	user, err = models.NewUser("admin", "admin@localhost", "password")
	if err != nil {
		log.Fatalln(err)
	}
	repos.Users.Insert(user)
	userAdmin = user.Id

	user, err = models.NewUser("user1", "user1@localhost", "password")
	if err != nil {
		log.Fatalln(err)
	}
	repos.Users.Insert(user)
	user1 = user.Id

	user, err = models.NewUser("user2", "user2@localhost", "password")
	if err != nil {
		log.Fatalln(err)
	}
	repos.Users.Insert(user)
	user2 = user.Id
}

func createPosts(conf config.Config) {
	createPost(user2, "user1", "text", "Hello World", "", "Lorem Ipsum", conf)
	createPost(user2, "user2", "link", "Stack Overflow", "http://stackoverflow.com/questions/9284144/google-go-vs-google-dart/", "", conf)
	createPost(user2, "user2", "link", "Slashdot", "http://slashdot.org/", "", conf)
	createPost(user2, "user2", "link", "GitHub", "https://github.com/dancannon/gorethink", "", conf)
	createPost(user2, "user2", "link", "BBC", "http://www.bbc.co.uk/news/world-europe-26677134", "", conf)
	createPost(user2, "user2", "link", "Vimeo", "http://vimeo.com/49718712", "", conf)
	createPost(user2, "user2", "link", "Twitter", "https://twitter.com/golang_news/status/447364460693684224", "", conf)
	createPost(user2, "user2", "link", "Flickr", "http://www.flickr.com/photos/photolupi/13309480293/in/explore-2014-03-21", "", conf)
	createPost(user2, "user2", "link", "Generic Blog", "http://blog.codinghorror.com/please-read-the-comments/", "", conf)
}

func createPost(authorID, author, typ, title, link, content string, conf config.Config) error {
	post := &models.Post{
		AuthorId:   authorID,
		AuthorName: author,
		Type:       typ,
		Title:      title,
		Link:       link,
		Content:    content,
		EmbedType:  "fetching",
		Created:    time.Now(),
		Modified:   time.Now(),
	}

	err := repos.Posts.Store(post)
	if err != nil {
		return err
	}

	// Start processing
	go lib.GeneratePostUrl(*post)
	go lib.ExtractLinkContent(conf.GoFetch, *post)

	return nil
}

func createRules() {
	r.Db("news").Table("rules").Insert(r.Json(`
		[{"AuthorId":"` + userAdmin + `","AuthorName":"admin","Created":{"$reql_type$":"TIME","epoch_time":1395517139,"timezone":"+00:00"},"Dislikes":0,"Host":"bbc.co.uk","Id":"da52cffa-f41c-4bce-a91f-8e999da5df41","Likes":0,"Modified":{"$reql_type$":"TIME","epoch_time":1395517139,"timezone":"+00:00"},"Name":"BBC News Article","PathPattern":"/news/.*","Type":"text","Url":"http://www.bbc.co.uk/news/world-europe-26677134","Values":[{"id":"selector","name":"title","params":{"attribute":"","restype":"first","selector":".story-header"},"type":"extractor"},{"id":"selector_text","name":"text","params":{"attribute":"","restype":"","selector":".introduction"},"type":"extractor"}],"id":"da52cffa-f41c-4bce-a91f-8e999da5df41"},{"AuthorId":"` + userAdmin + `","AuthorName":"admin","Created":{"$reql_type$":"TIME","epoch_time":1395517663,"timezone":"+00:00"},"Dislikes":0,"Host":"slashdot.org","Likes":0,"Modified":{"$reql_type$":"TIME","epoch_time":1395517663,"timezone":"+00:00"},"Name":"Slashdot Links","PathPattern":"/.*","Type":"general","Url":"http://slashdot.org","Values":[{"name":"title","type":"value","value":"Slashdot"},{"id":"selector","name":"content","params":{"attribute":"","restype":"merge","selector":"h2.story > span:nth-of-type(1) > a"},"type":"extractor"}],"id":"c0d5fb78-78d7-4752-927b-966ff8fea86c"},{"AuthorId":"` + userAdmin + `","AuthorName":"admin","Created":{"$reql_type$":"TIME","epoch_time":1395517663,"timezone":"+00:00"},"Dislikes":0,"Host":"beta.slashdot.org","Likes":0,"Modified":{"$reql_type$":"TIME","epoch_time":1395517663,"timezone":"+00:00"},"Name":"Beta Slashdot Links","PathPattern":"/.*","Type":"general","Url":"http://beta.slashdot.org","Values":[{"name":"title","type":"value","value":"Slashdot"},{"id":"selector","name":"content","params":{"attribute":"","restype":"merge","selector":".story-header > h1 > a"},"type":"extractor"}],"id":"c0d5fb78-78d7-4752-927b-966ff8fea86d"}]
	`)).RunWrite(infrastructure.RethinkDB())
}
