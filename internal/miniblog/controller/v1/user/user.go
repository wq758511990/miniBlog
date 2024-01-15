package user

import (
	"myMiniBlog/internal/miniblog/biz"
	"myMiniBlog/internal/miniblog/store"
	"myMiniBlog/pkg/auth"
	proto "myMiniBlog/pkg/proto/miniblog/v1"
)

type UserController struct {
	b biz.IBiz
	a *auth.Authz
	proto.UnimplementedMiniBlogServer
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{b: biz.NewBiz(ds), a: a}
}
