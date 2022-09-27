package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/trucktrace/pkg/logger"
)

const (
	signInKey = "truckTraceKey"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId    int    `json:"user_id"`
	CompanyId int    `json:"company_id"`
	Role      string `json:"user_role"`
}

func ParseToken(accessToken string) (tokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.ErrorLogger("ParseToken", "Can't parse token").Error("Token is invalid")
			return nil, errors.New("invalid token")
		}

		//logger.SystemLoggerInfo("ParseToken").Info("SignIn key was returned")
		return []byte(signInKey), nil
	})

	if err != nil {
		logger.SystemLoggerError("ParseToken", "Can't get signIn key").Error("Error - " + err.Error())
		return tokenClaims{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		logger.SystemLoggerError("ParseToken", "Token claims wrong type").Error("token claims are not of type *tokenClaims")
		return tokenClaims{}, errors.New("token claims are not of type *tokenClaims")
	}

	//logger.SystemLoggerInfo("ParseToken").Info("returned claims.UserId")

	return *claims, nil
}

func GetUserIdFromToken(ctx *gin.Context) (int, error) {

	value, isCookie := ctx.Request.Header["Set-Cookie"]

	if !isCookie {
		return 0, errors.New("user is not authorized - Cookie is empty")
	}

	claims, err := ParseToken(value[0])

	if err != nil {
		return 0, errors.New("token claims wrong type")
	}

	return claims.UserId, nil

}
