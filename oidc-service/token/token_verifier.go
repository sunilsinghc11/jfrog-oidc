package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

type GithubClaim struct {
	jwt.RegisteredClaims
}

func ParseToken(tokenValue string, getKeyFunc func(*jwt.Token) (interface{}, error)) (GithubClaim, error) {
	var githubClaim = GithubClaim{}
	_, err := jwt.ParseWithClaims(tokenValue, &githubClaim, getKeyFunc)
	if err != nil {
		log.Println("token: " + tokenValue)
		return GithubClaim{}, err
	}
	return githubClaim, nil
}

func VerifyToken(token string, getKeyFunc func(*jwt.Token) (interface{}, error)) (GithubClaim, error) {
	githubClaim, err := ParseToken(token, getKeyFunc)
	if err != nil {
		log.Println("failed to parse token err: " + err.Error())
		return GithubClaim{}, err
	}
	// verify audience
	var aud []byte
	if githubClaim.RegisteredClaims.Audience[0] != "access-oidc-poc" {
		log.Println("audience verification failed expected: access-oidc-poc got: " + string(aud))
		return GithubClaim{}, fmt.Errorf("audience verification failed expected")
	}
	// verify issuer

	return githubClaim, nil
}
