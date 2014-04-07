package repos

import (
	r "github.com/dancannon/gorethink"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type UserTokensRepo struct {
}

func NewUserTokensRepo() UserTokensRepo {
	return UserTokensRepo{}
}

func (repo *UserTokensRepo) FindToken(username string, secret string) (UserToken, error) {
	var token UserToken
	row, err := r.Table("user_tokens").GetAllByIndex("token", []interface{}{username, secret}).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return token, err
	}

	err = row.Scan(&token)

	return token, err
}

func (repo *UserTokensRepo) Insert(token *UserToken) error {
	response, err := r.Table("user_tokens").Insert(token).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if token.Id == "" && len(response.GeneratedKeys) == 1 {
		token.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *UserTokensRepo) Delete(id string) error {
	return r.Table("user_tokens").Get(id).Delete().Exec(infrastructure.RethinkDB())
}
