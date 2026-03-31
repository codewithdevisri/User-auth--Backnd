package dto

import (
	"fmt"
	"github.com/Ayan25844/netflix/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Payload struct {
	Name     string
	Email    string
	Password string
}

type JwtClaims struct {
	ID   string
	Name string
	Role []string
	jwt.StandardClaims
}

type LoginRequest struct {
	Email    string
	Password string
}

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(config.Ip, true) {
		return nil
	}
	return fmt.Errorf("token is invalid")
}
