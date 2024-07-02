package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main.go/src/configs/app"
	"main.go/src/shared/helpers"
	"main.go/src/shared/message/errors"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string, role ...string) gin.HandlerFunc {
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
				if len(role) > 0 && role != nil {
					flag := false
					var roles []string
					for _, r := range role {
						if strings.ToLower(claims.Role) == r || r == "" {
							flag = true
							break
						}
						roles = append(roles, r)
					}
					if !flag {
						c.JSON(http.StatusUnauthorized, errors.ErrorWrapper{Message: fmt.Sprintf("Only %v Can Access !", roles)})
						c.Abort()
						return
					}
				}
				c.Set("x-user-id", claims.ID)
				c.Set("x-user-role", claims.Role)
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
