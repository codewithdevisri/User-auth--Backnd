package token

import (
	"fmt"
	"github.com/Ayan25844/netflix/config"
	"github.com/Ayan25844/netflix/dto"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	jwtPrivateToken = "SecretTokenSecretToken"
	ip              = config.Ip
)

func GenerateToken(claims *dto.JwtClaims, expirationTime time.Time) (string, error) {
	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()
	claims.Issuer = ip
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtPrivateToken))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString, origin string) (bool, *dto.JwtClaims) {
	claims := &dto.JwtClaims{}
	token, _ := getTokenFromString(tokenString, claims)
	if token.Valid {
		if e := claims.Valid(); e == nil {
			return true, claims
		}
	}
	return false, claims
}

func getTokenFromString(tokenString string, claims *dto.JwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtPrivateToken), nil
	})
}
