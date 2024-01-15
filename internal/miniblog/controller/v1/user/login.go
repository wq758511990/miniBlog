package user

import (
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/log"
	v1 "myMiniBlog/pkg/api/miniblog/v1"
)

func (ctl *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	resp, err := ctl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)
}
