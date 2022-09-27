package middlewares

import (
	"context"
	"github.com/trucktrace/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/pkg/token"
)

func CheckSomePermission(privilege string, userRepo repository.UserRepository, groupRepo repository.GroupRepository) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		value, isCookie := ctx.Request.Header["Set-Cookie"]

		if !isCookie {
			res := response.BuildErrorResponse("User is not authorized", "Cookie is empty", response.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, err := token.ParseToken(value[0])

		if err != nil {
			res := response.BuildErrorResponse("Can't parse token", "Can't parse token", response.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userId := strconv.Itoa(claims.UserId)
		user, _ := userRepo.GetUserByIdForMiddleware(userId)

		if user.Role == "ADMIN" {
			ctx.Set("currentUser", user)
			c := context.Background()
			logger.MyUserContext = context.WithValue(c, logger.CtxKey{}, user)

			ctx.Next()
			return
		}

		var groups []models.Groups

		for _, value := range user.Groups {
			group, _ := groupRepo.GetGroupByIdForMiddleware(strconv.Itoa(int(value)))
			groups = append(groups, group)
		}

		for _, group := range groups {

			for _, val := range group.Priveleges {
				if val == privilege {
					ctx.Set("currentUser", user)

					c := context.Background()
					logger.MyUserContext = context.WithValue(c, logger.CtxKey{}, user)

					ctx.Next()
				}
			}
		}

		res := response.BuildErrorResponse("User doesn't have access", "User doesn't have access", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)

	}
}
