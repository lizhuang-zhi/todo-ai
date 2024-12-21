package middleware

import (
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"todo-ai/core/logger"

	"github.com/gin-gonic/gin"
)

// Recovery
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := r.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if !brokenPipe {
					logger.Errorf("Exception: %v\nStack:%v", r, string(debug.Stack()))
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(r.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
