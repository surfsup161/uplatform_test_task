package repo

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/surfsup161/uplatform_test_task/memcached"
	"github.com/surfsup161/uplatform_test_task/model"
)

type MemcachedStorage struct {
	m *memcached.Memcached
}

func NewMemcachedStorage(m *memcached.Memcached) Repo {
	return &MemcachedStorage{
		m: m,
	}
}

func (ms *MemcachedStorage) GetUser(id string) (*model.User, error) {
	data, err := ms.m.Get(id)
	if err != nil {
		return nil, err
	}

	userInfo := &model.User{}
	err = jsoniter.Unmarshal(data, userInfo)

	return userInfo, err
}

func (ms *MemcachedStorage) SetUser(user *model.User) error {
	data, err := jsoniter.Marshal(user)
	if err != nil {
		return err
	}

	return ms.m.Set(user.Id, data, 100)
}

func (ms *MemcachedStorage) DeleteUser(id string) error {
	return ms.m.Delete(id)
}
