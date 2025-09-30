package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chgrape/vaultpp/internal/repository"
	"github.com/chgrape/vaultpp/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

type LoginForm struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserValidator struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Role     string `json:"role" validate:"required,oneof=scribe member"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *UserService) Register(user UserValidator, ctx context.Context) (int, error) {
	err := validation.Instance().Struct(user)
	if err != nil {
		return 0, err
	}

	PassHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password")
	}

	id, err := s.Repo.CreateUser(repository.User{
		Name:     user.Name,
		Role:     user.Role,
		PassHash: string(PassHash)}, ctx)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *UserService) Login(user LoginForm, ctx context.Context) (string, error) {
	err := validation.Instance().Struct(user)
	if err != nil {
		return "", err
	}

	userData, err := s.Repo.GetUserPassHash(user.Name, ctx)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		UserID: strconv.Itoa(userData.ID),
		Role:   userData.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(userData.PassHash), []byte(user.Password)) != nil {
		return "", err
	}

	return tokenString, nil
}
