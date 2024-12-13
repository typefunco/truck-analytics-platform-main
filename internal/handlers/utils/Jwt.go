package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("5OGMnhB3g7<JeJzr+EesQ};0U_9sJOIO")

func CreateJWT(login string, password string) (string, error) {

	if login == "foton-trucks" || password == "foton1996" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"login": login,
				"exp":   time.Now().Add(time.Hour * 24 * 30 * 3).Unix(),
			})

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}

	return "", errors.New("Wrong password or login")
}

func VerifyJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
