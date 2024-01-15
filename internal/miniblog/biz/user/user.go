package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/marmotedu/miniblog/pkg/auth"
	"myMiniBlog/internal/miniblog/store"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/log"
	"myMiniBlog/internal/pkg/model"
	"myMiniBlog/pkg/api/miniblog/v1"
	"myMiniBlog/pkg/token"
	"regexp"
	"time"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ListUser(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds}
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}

func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}
	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}
	return nil
}

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, err
	}
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t}, nil
}

func (b *userBiz) ListUser(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}
	var ans []*v1.UserInfo
	for _, item := range list {
		ans = append(ans, &v1.UserInfo{
			Username:  item.Username,
			Nickname:  item.Nickname,
			Email:     item.Email,
			Phone:     item.Phone,
			PostCount: 5,
			CreatedAt: item.CreatedAt.Format(time.DateTime),
			UpdatedAt: item.UpdatedAt.Format(time.DateTime),
		})
	}
	return &v1.ListUserResponse{
		TotalCount: count,
		Users:      ans,
	}, nil

}
