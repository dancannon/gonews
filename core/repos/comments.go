package repos

import (
	r "github.com/dancannon/gorethink"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type CommentsRepo struct {
}

func NewCommentsRepo() CommentsRepo {
	return CommentsRepo{}
}

func (repo *CommentsRepo) FindAll() ([]*Comment, error) {
	var comments []*Comment
	rows, err := r.Table("comments").Run(infrastructure.RethinkDB())
	if err != nil {
		return comments, err
	}

	err = rows.ScanAll(&comments)
	if err != nil {
		return comments, err
	}

	return comments, err
}

func (repo *CommentsRepo) FindByPost(post string) ([]Comment, error) {
	var comments []Comment
	query := r.Table("comments").Filter(r.Row.Field("Post").Eq(post))
	query = query.Filter(r.Row.Field("Depth").Eq(0)).OrderBy(r.Desc(orderByTop))
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return comments, err
	}

	err = rows.ScanAll(&comments)
	if err != nil {
		return comments, err
	}

	return comments, err
}

func (repo *CommentsRepo) FindChildren(id string) ([]Comment, error) {
	var comments []Comment
	query := r.Table("comments").Filter(r.Row.Field("Parent").Eq(id))
	query = query.OrderBy(r.Desc(orderByTop))
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return comments, err
	}

	err = rows.ScanAll(&comments)
	if err != nil {
		return comments, err
	}

	return comments, err
}

func (repo *CommentsRepo) FindById(id string) (*Comment, error) {
	var comment = new(Comment)
	row, err := r.Table("comments").Get(id).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return comment, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(comment)

	return comment, err
}
func (repo *CommentsRepo) Store(comment *Comment) error {
	response, err := r.Table("comments").Insert(comment).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if comment.Id == "" && len(response.GeneratedKeys) == 1 {
		comment.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *CommentsRepo) Update(comment *Comment) error {
	_, err := r.Table("comments").Get(comment.Id).Update(comment).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepo) UpdateVote(id string, voteType string, undo, toggle bool) error {
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

	_, err := r.Table("comments").Get(id).Update(data).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepo) UpdateFields(id string, fields map[string]interface{}) error {
	_, err := r.Table("comments").Get(id).Update(fields).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepo) Delete(id string) error {
	return r.Table("comments").Get(id).Delete().Exec(infrastructure.RethinkDB())
}

func (repo *CommentsRepo) DeleteAll() error {
	return r.Table("comments").Delete().Exec(infrastructure.RethinkDB())
}
