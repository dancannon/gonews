package repos

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"regexp"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type RulesRepo struct {
}

func NewRulesRepo() RulesRepo {
	return RulesRepo{}
}

func (repo *RulesRepo) FindAll() ([]*Rule, error) {
	var rules []*Rule
	rows, err := r.Table("rules").Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindByHost(host string) ([]*Rule, error) {
	var rules []*Rule
	query := r.Table("rules").Filter(r.Row.Field("Host").Eq(host))

	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindByQueryHostAndPage(q string, host string, page int, count int) ([]*Rule, error) {
	var rules []*Rule
	query := r.Table("rules").Filter(r.Row.Field("Host").Eq(host)).Filter(func(row r.RqlTerm) r.RqlTerm {
		return row.Field("Name").Match(fmt.Sprintf("(?i)%s", regexp.QuoteMeta(q)))
	}).OrderBy(r.Desc(orderByTop)).Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindByHostAndPage(host string, page int, count int) ([]*Rule, error) {
	var rules []*Rule
	query := r.Table("rules").Filter(r.Row.Field("Host").Eq(host)).OrderBy(r.Desc(orderByTop))
	query = query.Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindByQueryAndPage(q string, page int, count int) ([]*Rule, error) {
	var rules []*Rule
	query := r.Table("rules").Filter(func(row r.RqlTerm) r.RqlTerm {
		return r.Expr([]string{"name", "host"}).Contains(func(key r.RqlTerm) r.RqlTerm {
			return row.Field(key).CoerceTo("string").Match(fmt.Sprintf("(?i)%s", regexp.QuoteMeta(q)))
		})
	}).OrderBy(r.Desc(orderByTop)).Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindByPage(page int, count int) ([]*Rule, error) {
	var rules []*Rule
	query := r.Table("rules").OrderBy(r.Desc(orderByTop))
	query = query.Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return rules, err
	}

	err = rows.ScanAll(&rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (repo *RulesRepo) FindById(id string) (*Rule, error) {
	var rule = new(Rule)
	row, err := r.Table("rules").Get(id).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return rule, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(rule)

	return rule, err
}
func (repo *RulesRepo) Store(rule *Rule) error {
	response, err := r.Table("rules").Insert(rule).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if rule.Id == "" && len(response.GeneratedKeys) == 1 {
		rule.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *RulesRepo) Update(rule *Rule) error {
	_, err := r.Table("rules").Get(rule.Id).Update(rule).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *RulesRepo) UpdateVote(id string, voteType string, undo, toggle bool) error {
	var data = make(map[string]interface{})
	if voteType == VoteTypeLike {
		if undo {
			data["Likes"] = r.Branch(
				r.Row.Field("Likes").Le(0),
				0,
				r.Row.Field("Likes").Sub(1),
			)
		} else {
			data["Likes"] = r.Row.Field("Likes").Add(1)
			if toggle {
				data["Dislikes"] = r.Branch(
					r.Row.Field("Dislikes").Le(0),
					0,
					r.Row.Field("Dislikes").Sub(1),
				)
			}
		}
	} else if voteType == VoteTypeDislike {
		if undo {
			data["Dislikes"] = r.Branch(
				r.Row.Field("Dislikes").Le(0),
				0,
				r.Row.Field("Dislikes").Sub(1),
			)
		} else {
			data["Dislikes"] = r.Row.Field("Dislikes").Add(1)
			if toggle {
				data["Likes"] = r.Branch(
					r.Row.Field("Likes").Le(0),
					0,
					r.Row.Field("Likes").Sub(1),
				)
			}
		}
	} else {
		return nil
	}

	_, err := r.Table("rules").Get(id).Update(data).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *RulesRepo) UpdateFields(id string, fields map[string]interface{}) error {
	_, err := r.Table("rules").Get(id).Update(fields).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *RulesRepo) Delete(id string) error {
	return r.Table("rules").Get(id).Delete().Exec(infrastructure.RethinkDB())
}

func (repo *RulesRepo) DeleteAll() error {
	return r.Table("rules").Delete().Exec(infrastructure.RethinkDB())
}
