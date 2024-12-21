package api

import (
	"net/http"
	"strings"
	"todo-ai/common"
	"todo-ai/common/auth"

	"github.com/gin-gonic/gin"
)

// 权限检查
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if common.Config.System.AuthOpen && !tokenCheck(c) {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"accessKey": common.Config.System.AuthAccessKey})
			return
		}
		c.Next()
	}
}

func tokenCheck(c *gin.Context) bool {
	params := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(params) < 2 {
		return false
	}
	token := params[1]

	if token == "null" {
		return false
	}

	// OpenAPI Token
	if token == common.Config.System.OpenAPIToken {
		c.Set(auth.CtxAccessToken, token)
		return true
	}

	return false
}
