package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"myMiniBlog/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		if requestID == "" {
			requestID = uuid.NewV4().String()
		}

		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(known.XRequestIDKey, requestID)

		// 将 RequestID 保存在 HTTP 返回头中，Header 的键为 `X-Request-ID`
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}
