package core

import (
	"github.com/gin-gonic/gin"
	"myMiniBlog/internal/pkg/errno"
	"net/http"
)

func CommonRsp(code string, message string) *Response {
	return &Response{Code: code, Message: message, Data: nil}
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) withData(data interface{}) *Response {
	return &Response{
		Code:    r.Code,
		Message: r.Message,
		Data:    data,
	}
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		hcode, code, message := errno.Decode(err)
		c.JSON(hcode, CommonRsp(code, message))
		return
	}
	c.JSON(http.StatusOK, CommonRsp("operate success", "操作成功").withData(data))
}
