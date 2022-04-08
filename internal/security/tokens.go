package security

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

func GetToken(tokenAuth *jwtauth.JWTAuth, id int64) (string, error) {

	var err error
	var token string

	if _, token, err = tokenAuth.Encode(map[string]interface{}{"id": id}); err != nil {
		fmt.Printf("Security :: GetToken error :: %v\n", err)
		return "", err
	}

	return token, nil

}

func GetIdFromContext(ctx context.Context) (int64, error) {

	var id float64
	var ok bool
	var err error
	var claims jwt.MapClaims

	if _, claims, err = jwtauth.FromContext(ctx); err != nil {
		fmt.Printf("GetIdFromContext error :: %v\n", err)
		return 0, err
	}

	if id, ok = claims["id"].(float64); !ok {
		err = errors.New("can not cast id claims to float64")
		fmt.Printf("GetIdFromContext error :: %v\n", err)
		return 0, err
	}

	return int64(id), nil

}

type userClaims struct {
	UserId int64
	jwt.StandardClaims
}

var key []byte = []byte(os.Getenv("CLIENT_KEY"))

func GetUserToken(userId int64) (string, error) {

	var err error
	var claims userClaims
	var token *jwt.Token
	var tokenString string

	claims = userClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 365 * 24 * time.Hour).Unix(),
			Issuer:    "maus",
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err = token.SignedString(key); err != nil {
		return "", err
	}

	return tokenString, nil

}

func ParseUserToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Signing method error")
		}
		return key, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims.UserId, nil
	}

	return 0, errors.New("Parsetoken failed")

}
