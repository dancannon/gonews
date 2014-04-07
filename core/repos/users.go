package repos

import (
	r "github.com/dancannon/gorethink"

	"github.com/dancannon/gonews/core/infrastructure"
	. "github.com/dancannon/gonews/core/models"
)

type UsersRepo struct {
}

func NewUsersRepo() UsersRepo {
	return UsersRepo{}
}

func (repo *UsersRepo) FindAll() ([]User, error) {
	var users []User
	rows, err := r.Table("users").Run(infrastructure.RethinkDB())
	if err != nil {
		return users, err
	}

	err = rows.ScanAll(&users)
	if err != nil {
		return users, err
	}

	return users, err
}

func (repo *UsersRepo) FindById(id string) (*User, error) {
	var user = new(User)
	row, err := r.Table("users").Get(id).RunRow(infrastructure.RethinkDB())
	if err != nil {
		return user, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(&user)

	return user, err
}

func (repo *UsersRepo) FindByUsername(username string) (*User, error) {
	var user = new(User)
	query := r.Table("users").GetAllByIndex("Username", username)

	row, err := query.RunRow(infrastructure.RethinkDB())
	if err != nil {
		return user, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(user)

	return user, err
}

func (repo *UsersRepo) FindByEmail(email string) (*User, error) {
	var user = new(User)
	query := r.Table("users").GetAllByIndex("Email", email)

	row, err := query.RunRow(infrastructure.RethinkDB())
	if err != nil {
		return user, err
	}

	if row.IsNil() {
		return nil, nil
	}

	err = row.Scan(user)

	return user, err
}

func (repo *UsersRepo) Insert(user *User) error {
	response, err := r.Table("users").Insert(user).RunWrite(infrastructure.RethinkDB())

	if err != nil {
		return err
	}

	// Find new ID of product if needed
	if user.Id == "" && len(response.GeneratedKeys) == 1 {
		user.Id = response.GeneratedKeys[0]
	}

	return nil
}

func (repo *UsersRepo) Delete(id string) error {
	return r.Table("users").Get(id).Delete().Exec(infrastructure.RethinkDB())
}

func (repo *UsersRepo) DeleteAll() error {
	return r.Table("users").Delete().Exec(infrastructure.RethinkDB())
}
