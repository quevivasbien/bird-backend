package api

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/quevivasbien/bird-backend/db"
)

var JWT_SECRET = []byte(os.Getenv("BIRD_JWT_SECRET"))

const JWT_EXPIRE_HOURS = 12

func getToken(user db.User) (string, time.Time, error) {
	expireTime := time.Now().Add(time.Hour * JWT_EXPIRE_HOURS)
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   user.Name,
			"admin": user.Admin,
			"exp":   expireTime.Unix(),
		},
	)
	token, err := claims.SignedString(JWT_SECRET)
	return token, expireTime, err
}
