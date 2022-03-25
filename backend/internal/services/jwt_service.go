package services

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// secret key being used to sign tokens
type JwtService struct {
	SecretKey []byte
}

type UserClaim struct {
	*jwt.StandardClaims
	Username string
}

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func (s *JwtService) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &UserClaim{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
		Username: username,
	}

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func (s *JwtService) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})
	if claims, ok := token.Claims.(UserClaim); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return "", err
		}

		username := claims.Username
		return username, nil
	} else {
		return "", err
	}
}
