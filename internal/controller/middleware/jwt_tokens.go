package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ellofae/authentication-deanery/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenData struct {
	Expiry   int64
	IssuedAt int64
	UserID   string
	State    string
}

var jwtSecretKey string

func InitJWTSecretKey(secretKey string) {
	jwtSecretKey = secretKey
}

func GenerateAccessToken(record_code int) (string, error) {
	logger := logger.GetLogger()

	current_time := time.Now()
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"record_code": strconv.Itoa(record_code),
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

func ParseToken(tokenString string) (*AccessTokenData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := claims["expiry"].(float64)
		issued_at := claims["issued_at"].(float64)

		return &AccessTokenData{
			Expiry:   int64(expiry),
			IssuedAt: int64(issued_at),
			UserID:   claims["record_code"].(string),
			State:    claims["state"].(string),
		}, nil
	} else {
		return nil, err
	}
}
