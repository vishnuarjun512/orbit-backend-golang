package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
}

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

func (j *JWTService) CreateToken(userID string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}

func (j *JWTService) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("Unexpected signing method")
			}

			return j.secret, nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	if claims.UserID == "" {
		return "", errors.New("User ID missing from token")
	}

	return claims.UserID, nil
}
