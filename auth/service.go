package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Service interface {
	VerifyToken(c *gin.Context) (interface{}, error)
	GenerateToken2(user interface{}, access_token string) string
	GenerateToken(id int, email string, username string, is_organizer int) string
}

func VerifyToken(stringToken string) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(os.Getenv("AUTHSECRETKEY")), nil
	})
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResponse
	}
	return token.Claims.(jwt.MapClaims), nil
}

func GenerateToken2(user interface{}, access_token string) string {
	expirationTime := time.Now().Add(1 * time.Minute).Unix()
	claims := jwt.MapClaims{
		"user":         user,
		"access_token": access_token,
		"expired_at":   expirationTime,
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singnedToken, _ := parseToken.SignedString([]byte(os.Getenv("AUTHSECRETKEY")))
	return singnedToken
}

func GenerateToken(id int, email string, username string, is_organizer int) string {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id":            id,
		"email":              email,
		"username":           username,
		"event_is_organizer": is_organizer,
		"expired_at":         expirationTime,
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singnedToken, _ := parseToken.SignedString([]byte(os.Getenv("AUTHSECRETKEY")))
	return singnedToken
}