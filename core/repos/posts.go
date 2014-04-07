package repos

import (
	r "github.com/dancannon/gorethink"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type PostsRepo struct {
}

func NewPostsRepo() PostsRepo {
	return PostsRepo{}
}

func (repo *PostsRepo) FindAll() ([]*Post, error) {
	var posts []*Post
	rows, err := r.Table("posts").Run(infrastructure.RethinkDB())
	if err != nil {
		return posts, err
	}

	err = rows.ScanAll(&posts)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (repo *PostsRepo) FindById(id string) (*Post, error) {
	var post = new(Post)
	row, err := r.Table("posts").Get(id).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return post, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(post)

	return post, err
}

func (repo *PostsRepo) FindByUrl(url string) (*Post, error) {
	var post = new(Post)
	row, err := r.Table("posts").Filter(r.Row.Field("Url").Eq(url)).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return post, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(post)

	return post, err
}

func (repo *PostsRepo) FindPopularByPage(page int, count int) ([]*Post, error) {
	var posts []*Post
	query := r.Table("posts").OrderBy(r.Desc(orderByPopular))
	query = query.Map(func(row r.RqlTerm) r.RqlTerm {
		return row.Merge(map[string]interface{}{
			"Author": r.Table("users").Get(row.Field("Author")).Field("Username").Default(""),
		})
	})
	query = query.Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return posts, err
	}

	err = rows.ScanAll(&posts)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (repo *PostsRepo) FindTopByPage(page int, count int) ([]*Post, error) {
	var posts []*Post
	query := r.Table("posts").OrderBy(r.Desc(orderByTop))
	query = query.Map(func(row r.RqlTerm) r.RqlTerm {
		return row.Merge(map[string]interface{}{
			"Author": r.Table("users").Get(row.Field("Author")).Field("Username").Default(""),
		})
	})
	query = query.Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return posts, err
	}

	err = rows.ScanAll(&posts)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (repo *PostsRepo) FindNewByPage(page int, count int) ([]*Post, error) {
	var posts []*Post
	query := r.Table("posts").OrderBy(r.Desc(orderByNew))
	query = query.Skip((page - 1) * count).Limit(count)
	rows, err := query.Run(infrastructure.RethinkDB())
	if err != nil {
		return posts, err
	}

	err = rows.ScanAll(&posts)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (repo *PostsRepo) Store(post *Post) error {
	response, err := r.Table("posts").Insert(post).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if post.Id == "" && len(response.GeneratedKeys) == 1 {
		post.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *PostsRepo) Update(post *Post) error {
	_, err := r.Table("posts").Get(post.Id).Update(post).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostsRepo) UpdateVote(id string, voteType string, undo, toggle bool) error {
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

	_, err := r.Table("posts").Get(id).Update(data).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostsRepo) UpdateFields(id string, fields map[string]interface{}) error {
	_, err := r.Table("posts").Get(id).Update(fields).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostsRepo) Delete(id string) error {
	return r.Table("posts").Get(id).Delete().Exec(infrastructure.RethinkDB())
}

func (repo *PostsRepo) DeleteAll() error {
	return r.Table("posts").Delete().Exec(infrastructure.RethinkDB())
}
