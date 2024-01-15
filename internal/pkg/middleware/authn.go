package middleware

import (
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/known"
	"myMiniBlog/pkg/token"
)

func Authn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := token.ParseRequest(ctx)
		if err != nil {
			core.WriteResponse(ctx, errno.ErrTokenInvalid, nil)
			ctx.Abort()
			return
		}
		ctx.Set(known.XUsernameKey, username)
		ctx.Next()
	}
}
