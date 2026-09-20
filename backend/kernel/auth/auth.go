package authjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"uid"`
	Role      string `json:"role"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

type TokenMaker struct {
	secret []byte
	ttl    time.Duration
}

func NewHS256(secret string, ttl time.Duration) (*TokenMaker, error) {
	if len(secret) < 32 {
		return nil, jwt.ErrInvalidKey
	}
	return &TokenMaker{secret: []byte(secret), ttl: ttl}, nil
}

func (t *TokenMaker) Issue(userID, role string) (string, error) {
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID, Role: role, TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
		},
	})
	return tok.SignedString(t.secret)
}

func (t *TokenMaker) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(rt *jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
