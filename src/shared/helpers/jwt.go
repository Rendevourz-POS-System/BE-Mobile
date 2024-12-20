package helpers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	appConfig "main.go/src/configs/app"
	User "main.go/src/domains/user/entities"
	"time"
)

func generateTokenTime(expiry int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry)))
}

func generateTokenTimeMinute(expiry int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry)))
}

func GenerateRefreshToken(user *User.User, expiry int) (string, error) {
	claimsRefresh := &User.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: generateTokenTime(expiry),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	tokenString, err := token.SignedString([]byte(appConfig.GetConfig().AccessToken.AccessTokenSecret))
	if err != nil {
		return "", errors.New("error when generate token")
	}
	return tokenString, nil
}

func GenerateToken(user *User.User) (string, error) {
	expiry := generateTokenTime(appConfig.GetConfig().AccessToken.AccessTokenExpireHour)
	claims := &User.JwtCustomClaims{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		Role:     user.Role,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiry,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(appConfig.GetConfig().AccessToken.AccessTokenSecret))
	if err != nil {
		return "", errors.New("error when generate token")
	}
	return tokenString, nil
}

func VerifyToken(requestToken, secret string) (bool, error) {
	errResponse := errors.New("middleware Error: failed to authenticate token ! ")
	_, err := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, errResponse
	}
	return true, nil
}

func ClaimsTokenData(requestToken, secret string) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	token, _ := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
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

func ExtractIDFromToken(requestToken, secret string) (User.JwtCustomClaims, error) {
	errResponse := errors.New("sign in to proceed")
	token, _ := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	data := User.JwtCustomClaims{}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return data, errResponse
	}
	data = User.JwtCustomClaims{
		ID:    claims["Id"].(string),
		Email: claims["Email"].(string),
		Otp:   nil,
		Role:  claims["Role"].(string),
	}
	return data, nil
}

func GenerateJwtTokenForVerificationEmail(Otp *int, id, email, secretCode string) (string, error) {
	claims := &User.JwtEmailClaims{
		ID:    id,
		Email: email,
		Otp:   Otp,
		Nonce: secretCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: generateTokenTimeMinute(appConfig.GetConfig().AccessToken.VerificationTokenExpireHour),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(appConfig.GetConfig().AccessToken.AccessTokenSecret))
	if err != nil {
		return "", errors.New("error when generate token")
	}
	return tokenString, nil
}

func ClaimsJwtTokenForVerificationEmail(requestToken string) (*User.JwtEmailClaims, error) {
	errResponse := errors.New("sign in to proceed")
	token, _ := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(appConfig.GetConfig().AccessToken.AccessTokenSecret), nil
	})
	//fmt.Println("Masuk : ", token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		//fmt.Println("Error : ", errResponse)
		return nil, errResponse
	}
	//fmt.Println("Claims : ", claims)
	exps, _ := claims.GetExpirationTime()
	if time.Now().UTC().After(time.Unix(exps.Unix(), 0)) {
		return nil, errors.New("token has expired")
	}
	return &User.JwtEmailClaims{
		ID:    claims["Id"].(string),
		Email: claims["Email"].(string),
		Nonce: claims["Nonce"].(string),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: exps,
		},
	}, nil
}
