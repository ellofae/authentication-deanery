package middleware

import (
	"errors"
	"strconv"
	"time"

	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenData struct {
	Expiry     int64
	IssuedAt   int64
	RecordCode string
	Role       string
	State      string
}

var jwtSecretKey string

func InitJWTSecretKey(secretKey string) {
	jwtSecretKey = secretKey
}

func GenerateAccessToken(record_code int, role string) (string, error) {
	logger := logger.GetLogger()

	current_time := time.Now()
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"record_code": strconv.Itoa(record_code),
			"role":        role,
			"issued_at":   current_time.Unix(),
			"expiry":      current_time.Add(time.Minute * 30).Unix(),
			"state":       "access_token",
		})

	access_token, err := jwt_token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		logger.Print("Unable to generate an access token")
		return "", errors.New("unable to generate an access token")
	}
	return access_token, nil
}
