package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

const (
	key                     = "secret"
	refreshTokenExpireTime  = 3 * 24 * time.Hour // Refresh token expiration time
	accessTokenExpireTime   = 24 * time.Hour     // Access token expiration time
	accessTokenHeaderName   = "Authorization"    // Header name for sending the access token
	accessTokenHeaderPrefix = "Bearer "          // Prefix for the access token header value
	//refreshTokenCookieName  = "refresh_token"       // Name of the refresh token cookie
	//accessTokenCookieName   = "access_token"        // Name of the access token cookie
	//refreshTokenCookiePath  = "/"                   // Path for the refresh token cookie
	//accessTokenCookiePath   = "/"                   // Path for the access token cookie
	//refreshTokenHeaderName  = "X-Refresh-Token"     // Header name for sending the refresh token
)

func RefreshToken(id uint, tokenString string) (string, error) {
	expiry := time.Now().Add(refreshTokenExpireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId": id,
		"Exp":    expiry,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("error when generate token")
	}
	return tokenString, nil
}

func GenerateRefreshToken(id uint, token string) (string, error) {
	expiry := time.Now().Add(refreshTokenExpireTime).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":  id,
		"Token":   token,
		"Expired": expiry,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("error when generate refresh token")
	}
	return refreshTokenString, nil
}

func GenerateToken(id, email, roleName, divisionName string) (string, error) {
	expiry := time.Now().Add(accessTokenExpireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":   id,
		"Email":    email,
		"Role":     roleName,
		"Division": divisionName,
		"Expired":  expiry,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("error when generate token")
	}
	return tokenString, nil
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	headerToken := c.Request.Header.Get(accessTokenHeaderName)
	bearer := strings.HasPrefix(headerToken, accessTokenHeaderPrefix)
	if headerToken == "" || !bearer {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  errResponse.Error(),
			"errors": headerToken,
		})
		return nil, errResponse
	}
	if !bearer {
		return nil, errResponse
	}
	stringToken := strings.Split(headerToken, " ")[1]
	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(key), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errResponse
	}
	ExpiredAt := int64(claims["Expired"].(float64))
	if time.Now().UTC().After(time.Unix(ExpiredAt, 0)) {
		return nil, errors.New("token has expired")
	}
	return token.Claims.(jwt.MapClaims), nil
}

//func ParseToken(tokenString string) (*request.UserTokenRequest, error) {
//	claims := jwt.MapClaims{}
//	if strings.HasPrefix(tokenString, accessTokenHeaderPrefix) {
//		tokenString = strings.Split(tokenString, " ")[1]
//	}
//	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte(key), nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	if !token.Valid {
//		return nil, errors.New("invalid token")
//	}
//	// do something with decoded claims
//	for key, val := range claims {
//		fmt.Printf("Key: %v, value: %v\n", key, val)
//	}
//	return &request.UserTokenRequest{
//		UserId:       uint(claims["UserId"].(float64)),
//		Email:        claims["Email"].(string),
//		RoleName:     claims["Role"].(string),
//		DivisionName: claims["Division"].(string),
//		ExpiredAt:    int64(claims["Expired"].(float64)),
//		AccessToken:  tokenString,
//	}, nil
//}

//func GetTokenData(c *gin.Context) (*request.UserTokenRequest, error) {
//	claims := jwt.MapClaims{}
//	tokenString := c.Request.Header.Get(accessTokenHeaderName)
//	if strings.HasPrefix(tokenString, accessTokenHeaderPrefix) {
//		tokenString = strings.Split(tokenString, " ")[1]
//	}
//	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte(key), nil
//	})
//	if !token.Valid {
//		return nil, errors.New("invalid token")
//	}
//	return &request.UserTokenRequest{
//		UserId:       uint(claims["UserId"].(float64)),
//		Email:        claims["Email"].(string),
//		RoleName:     claims["Role"].(string),
//		DivisionName: claims["Division"].(string),
//		ExpiredAt:    int64(claims["Expired"].(float64)),
//		AccessToken:  tokenString,
//	}, nil
//}
