package repos

import (
	r "github.com/dancannon/gorethink"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type VotesRepo struct {
}

func NewVotesRepo() VotesRepo {
	return VotesRepo{}
}

func (repo *VotesRepo) FindAll() ([]*Vote, error) {
	var votes []*Vote
	rows, err := r.Table("votes").Run(infrastructure.RethinkDB())
	if err != nil {
		return votes, err
	}

	err = rows.ScanAll(&votes)
	if err != nil {
		return votes, err
	}

	return votes, err
}

func (repo *VotesRepo) FindByEntityAndUser(entityId, userId string) (*Vote, error) {
	var vote = new(Vote)
	row, err := r.Table("votes").Filter(
		r.Row.Field("Entity").Eq(entityId).And(r.Row.Field("User").Eq(userId)),
	).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return vote, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(vote)

	return vote, err
}

func (repo *VotesRepo) FindById(id string) (*Vote, error) {
	var vote = new(Vote)
	row, err := r.Table("votes").Get(id).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return vote, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(vote)

	return vote, err
}
func (repo *VotesRepo) Store(vote *Vote) error {
	response, err := r.Table("votes").Insert(vote).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if vote.Id == "" && len(response.GeneratedKeys) == 1 {
		vote.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *VotesRepo) Update(vote *Vote) error {
	_, err := r.Table("votes").Get(vote.Id).Update(vote).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *VotesRepo) UpdateFields(id string, fields map[string]interface{}) error {
	_, err := r.Table("votes").Get(id).Update(fields).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *VotesRepo) Delete(id string) error {
	return r.Table("votes").Get(id).Delete().Exec(infrastructure.RethinkDB())
}

func (repo *VotesRepo) DeleteAll() error {
	return r.Table("votes").Delete().Exec(infrastructure.RethinkDB())
}
