package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	todo "github.com/Tim-Masuda/rest_todo"
	"github.com/Tim-Masuda/rest_todo/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "tmtmtmtmtmtmtmtmtmmttmtmmtmt77"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour // time for token
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password) // хеширование пароля и запись его обратно
	return s.repo.CreateUser(user)
}

type tokenClaims struct { // база токена чтоб удобно обращиться было
	jwt.StandardClaims
	UserId int `json:"user_id"` // хран id \
}

// создание токена
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // через 12 часов престанет работать токен
			IssuedAt:  time.Now().Unix(),               // дата генерации токена
		},
		user.Id, // привязка токена по id
	})

	return token.SignedString([]byte(signingKey))
}

// парсинг токена 
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil // id пользотеля 
}

func generatePasswordHash(password string) string { // хешированя пароля
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
