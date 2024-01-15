package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/core"
	"myMiniBlog/internal/pkg/errno"
	"myMiniBlog/internal/pkg/log"
	v1 "myMiniBlog/pkg/api/miniblog/v1"
)

const defaultMethods = "(GET)|(POST)|(PUT)|(DELETE)"

func (ctl *UserController) Create(c *gin.Context) {
	log.C(c).Infow("Create user function called")
	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	if err := ctl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	if _, err := ctl.a.AddNamedPolicy("p", r.Username, "/v1/users/"+r.Username, defaultMethods); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
