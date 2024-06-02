package middlewares

import (
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(app.GetConfig().AccessToken.AccessTokenHeaderName)
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := helpers.VerifyToken(authToken, secret)
			if authorized {
				claims, err := helpers.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, errors.ErrorWrapper{Message: err.Error()})
					c.Abort()
					return
				}
				if len(role) > 0 && role != "" {
					if strings.ToLower(claims.Role) != role {
						c.JSON(http.StatusUnauthorized, errors.ErrorWrapper{Message: "Only " + role + " Can Access !"})
						c.Abort()
						return
					}
				}
				c.Set("x-user-id", claims.ID)
				c.Next()
				return
			}
			//fmt.Println("Auth token --> ", authorized)
			//fmt.Println("User ID --> ", userID)
			c.JSON(http.StatusUnauthorized, errors.ErrorWrapper{Message: err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, errors.ErrorWrapper{Message: "Not authorized"})
		c.Abort()
	}
}
