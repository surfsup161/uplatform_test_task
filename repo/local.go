package repo

import (
	"sync"

	"github.com/surfsup161/uplatform_test_task/model"
)

type LocalStorage struct {
	mu       sync.RWMutex
	usersMap map[string]*model.User
}

func NewLocalStorageRepo() Repo {
	return &LocalStorage{
		usersMap: make(map[string]*model.User),
	}
}

func (ls *LocalStorage) GetUser(id string) (*model.User, error) {
	ls.mu.RLock()
	userInfo := ls.usersMap[id]
	ls.mu.RUnlock()

	return userInfo, nil
}

func (ls *LocalStorage) SetUser(user *model.User) error {
	ls.mu.Lock()
	ls.usersMap[user.GetId()] = user
	ls.mu.Unlock()

	return nil
}

func (ls *LocalStorage) DeleteUser(id string) error {
	ls.mu.Lock()
	delete(ls.usersMap, id)
	ls.mu.Unlock()

	return nil
}
