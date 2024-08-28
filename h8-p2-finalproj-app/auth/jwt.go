package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtAppClaims struct {
	UserID uint
	jwt.RegisteredClaims
}
