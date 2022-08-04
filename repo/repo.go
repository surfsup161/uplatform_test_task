package repo

import "github.com/surfsup161/uplatform_test_task/model"

type Repo interface {
	GetUser(id string) (*model.User, error)
	SetUser(user *model.User) error
	DeleteUser(id string) error
}
