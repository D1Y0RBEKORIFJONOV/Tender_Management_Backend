package token

import (
	"awesomeProject/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Token()), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return claims, nil
}

func GetIdFromToken(r *http.Request) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["uid"]), 0
}

func GetEmailFromToken(r *http.Request) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["email"]), 0
}
