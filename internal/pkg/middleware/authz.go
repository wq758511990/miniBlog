package middleware

import (
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/known"
	"myMiniBlog/internal/pkg/log"
)

type Auther interface {
	Authorize(sub, obj, act string) (bool, error)
}

func Authz(a Auther) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sub := ctx.GetString(known.XUsernameKey)
		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		log.Debugw("Build authorize context", "sub", sub, "obj", obj, "act", act)
		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
	}
}
