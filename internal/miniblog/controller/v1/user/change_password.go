package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/log"
	v1 "myMiniBlog/pkg/api/miniblog/v1"
)

func (ctl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("Change password function called")
	var r v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctl.b.Users().ChangePassword(c, c.Param("name"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}
	core.WriteResponse(c, nil, nil)
}
