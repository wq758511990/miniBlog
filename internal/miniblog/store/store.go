package store

import (
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once
	S    *datastore
)

type IStore interface {
	DB() *gorm.DB
	Users() UserStore
}

type datastore struct {
	db *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db}
	})
	return S
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}
