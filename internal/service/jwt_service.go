package service

import (
	"fmt"
	"time"
	"card_manage/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT generation and validation.
type JWTService struct {
	secretKey      string
	expireDuration time.Duration
}

// CustomClaims are our custom claims, which includes standard claims and user-specific data.
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, expireDurationStr string) (*JWTService, error) {
	expireDuration, err := time.ParseDuration(expireDurationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token expire duration: %w", err)
	}

	return &JWTService{
		secretKey:      secretKey,
		expireDuration: expireDuration,
	}, nil
}

// GenerateToken generates a new JWT for a given user.
func (s *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "card_manage_platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates the given JWT string.
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
