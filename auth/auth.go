package auth

import (
	"fasnap-server-go/errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var mySigningKey = []byte("201801231158")

type myClaims struct {
	user string
	jwt.StandardClaims
}

func getSignedTokenStr(user string, expiresAt int64) (string, error) {
	standardClaims := jwt.StandardClaims{ExpiresAt: expiresAt}
	claims := myClaims{user, standardClaims}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func Authorize() gin.HandlerFunc {
	duration := time.Hour * 72

	return func(ctx *gin.Context) {
		user, password := ctx.Param("user"), ctx.Param("password")
		if user == "admin" && password == "admin" {
			tokenStr, err := getSignedTokenStr(user, time.Now().Add(duration).Unix())

			if err == nil {
				ctx.JSON(200, gin.H{"token": tokenStr})
				return
			}
		}

		ctx.AbortWithStatusJSON(200, errors.IncorrectUserPwd())
	}
}

func Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if len(tokenStr) == 0 {
			ctx.AbortWithStatusJSON(200, errors.CanntFindToken())
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if token != nil && token.Valid {
			ctx.Set("user", "admin")
			ctx.Next()
			return
		}

		if vErr, ok := err.(*jwt.ValidationError); ok {
			if vErr.Errors&jwt.ValidationErrorMalformed != 0 {
				ctx.AbortWithStatusJSON(200, errors.NotEvenToken())
				return
			} else if vErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				ctx.AbortWithStatusJSON(200, errors.ExpiredOrNotActiveToken())
				return
			}
		}
		ctx.AbortWithStatusJSON(200, errors.CanntHandleToken())
	}
}
