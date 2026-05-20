package middleware

import (
	"gin_blog_upload/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.APIResponse{
			Code: 7,
			Msg:  "服务器异常",
			Data: map[string]interface{}{},
		})
	})
}

