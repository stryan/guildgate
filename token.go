package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenClaim struct {
	Sponsor string `json:"sponsor_username"`
	jwt.StandardClaims
}

func generateToken(sponsor string) (string, error) {
	claim := tokenClaim{
		Sponsor: sponsor,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Unix() + 86400,
			Issuer:    "GuildGate",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(Conf.Secret))
	if err != nil {
		return "", err
	} else {
		return signedToken, nil
	}
}

func validateToken(tok string) (string, error) {
	token, err := jwt.ParseWithClaims(
		strings.TrimSpace(tok),
		&tokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Conf.Secret), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*tokenClaim)
	if !ok {
		return "", errors.New("Invalid token sponsor passed")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return "", errors.New("Token has expired")
	}
	log.Printf("Valid token received; sponsored by %v\n", claims.Sponsor)
	return claims.Sponsor, nil
}
