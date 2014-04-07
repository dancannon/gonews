package repos

var (
	Posts      PostsRepo
	Users      UsersRepo
	UserTokens UserTokensRepo
	Votes      VotesRepo
	Comments   CommentsRepo
	Rules      RulesRepo
)

func InitRepos() {
	Posts = NewPostsRepo()
	Users = NewUsersRepo()
	UserTokens = NewUserTokensRepo()
	Votes = NewVotesRepo()
	Comments = NewCommentsRepo()
	Rules = NewRulesRepo()
}
