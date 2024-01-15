package miniblog

import (
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/miniblog/controller/v1/user"
	"myMiniBlog/internal/miniblog/store"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/log"
	"myMiniBlog/internal/pkg/middleware"
	"myMiniBlog/pkg/auth"
)

func installRouters(g *gin.Engine) error {
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("healthz function called")
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})
	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	uc := user.New(store.S, authz)
	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{
		userV1 := v1.Group("/users")
		{
			userV1.POST("", uc.Create)
			userV1.PUT(":name/change-password", middleware.Authn(), middleware.Authz(authz), uc.ChangePassword)
			userV1.POST("/list", middleware.Authn(), uc.List)
		}
	}
	return nil
}
