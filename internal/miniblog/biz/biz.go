package biz

import (
	"myMiniBlog/internal/miniblog/biz/user"
	"myMiniBlog/internal/miniblog/store"
)

type IBiz interface {
	Users() user.UserBiz
}

var _ IBiz = (*biz)(nil)

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
