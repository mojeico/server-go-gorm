package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/response"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {

	return func(context *gin.Context) {

		_, isCookie := context.Cookie("accessToken")

		if isCookie != nil {
			res := response.BuildErrorResponse("User is not authorized", "Cookie is empty", response.EmptyObj{})
			context.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		context.Next()
	}
}
