package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/pkg/logger"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			logger.SystemLoggerError("CORSMiddleware", "user is not authorized").Error("Error - empty cookie")
			res := response.BuildErrorResponse("User is not authorized", "Cookie is empty", response.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusNoContent, res)
			return
		}
		c.Next()
	}
}
