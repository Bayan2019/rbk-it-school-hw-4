package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// 2. JWT

var (
	ErrTokenExpired  = errors.New("token is expired")
	ErrTokenNotFound = errors.New("token not found")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid claims")
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret []byte) *JWTManager {
	return &JWTManager{secret: secret}
}

// Payload должен содержать: { user_id, email, role, exp }
type Claims struct {
	UserID int64        `json:"user_id"`
	Email  string       `json:"email"`
	Role   domain.Roles `json:"role"`
	Exp    time.Time    `json:"exp"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID int64, email string, role domain.Roles) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Exp:    time.Now().Add(30 * time.Minute),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(90 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jwt-auth-lesson",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Validate(tokenString string) (*Claims, error) {
	// 6. Безопасность
	// - проверка подписи
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Best practice: never trust the algorithm from token blindly.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return claims, ErrInvalidClaims
	}

	// 6. Безопасность
	// - проверка exp
	if claims.Exp.After(time.Now().Add(30 * time.Minute)) {
		return claims, ErrTokenExpired
	}

	return claims, nil
}
