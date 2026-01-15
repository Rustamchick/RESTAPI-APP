package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"restapi-app"
	"time"

	"restapi-app/pkg/repository"

	"github.com/golang-jwt/jwt/v5"
)

const (
	salt       = "dw14312r3r2eqdf42tfewf4y6yw"
	signingKey = "fdhsufhdshu432ufed8wfu83@$#"
	TokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repos repository.Authorization
}

// create new service for authentification
func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos}
}

// use this method for signing up
func (s *AuthService) CreateUser(user restapi.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}

// use this method for signing in
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repos.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return " ", err
	}

	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
		UserId: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, errors.New("error parsing token")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
