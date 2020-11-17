package auth

import (
	"errors"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type JWTAuth struct {
	JWTToken string `json:"jwt_token"`
}

type ClaimsWithUID struct {
	jwtgo.StandardClaims
	TokenType string `json:"token_type"`
	UID       int    `json:"uid"`
}

func ValidateGetUIDJWT(secretkey, jwt string) (int, error) {
	token, err := jwtgo.ParseWithClaims(jwt, &ClaimsWithUID{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	fmt.Println(secretkey)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*ClaimsWithUID); ok && token.Valid {
		fmt.Println(claims.TokenType, claims.UID)
		return claims.UID, nil
	} else {
		return 0, errors.New("Wrong claims structure")
	}
}
