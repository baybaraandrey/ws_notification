package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	jwtgo "github.com/dgrijalva/jwt-go"
)

var (
	tokenString = flag.String("token", "", "-token jwt token")
)

type CustomClaimsExample struct {
	jwtgo.StandardClaims
	TokenType string `json:"token_type"`
	UID       int    `json:"uid"`
}

func main() {
	flag.Parse()

	// Parse the token
	token, err := jwtgo.ParseWithClaims(*tokenString, &CustomClaimsExample{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return []byte("lSu++8iV5rwcgNEKoQw3fTp4GFHUq1S/lTbiWB2TUEo="), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(*CustomClaimsExample); ok && token.Valid {
		fmt.Println(claims.TokenType, claims.UID)
	}
}
