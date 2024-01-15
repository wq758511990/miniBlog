package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/log"
	v1 "myMiniBlog/pkg/api/miniblog/v1"
	proto "myMiniBlog/pkg/proto/miniblog/v1"
	"time"
)

func (ctl *UserController) List(ctx *gin.Context) {
	log.C(ctx).Infow("List user function called")

	var r v1.ListUserRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	resp, err := ctl.b.Users().ListUser(ctx, r.Offset, r.Limit)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, resp)
}

// ListUser 返回用户列表，只有 root 用户才能获取用户列表.
func (ctrl *UserController) ListUser(ctx context.Context, r *proto.ListUserRequest) (*proto.ListUserResponse, error) {
	log.C(ctx).Infow("ListUser function called")

	resp, err := ctrl.b.Users().ListUser(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}

	var users []*proto.UserInfo
	for _, user := range resp.Users {
		createdAt, _ := time.Parse(time.DateTime, user.CreatedAt)
		updatedAt, _ := time.Parse(time.DateTime, user.UpdatedAt)
		users = append(users, &proto.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			PostCount: user.PostCount,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	return &proto.ListUserResponse{
		TotalCount: resp.TotalCount,
		Users:      users,
	}, nil
}
