package store

import (
	"context"
	"gorm.io/gorm"
	"myMiniBlog/internal/pkg/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)
}

type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

func (u *users) List(ctx context.Context, offset, limit int) (count int64, res []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&res).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}
