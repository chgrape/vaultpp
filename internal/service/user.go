package service

import (
	"context"
	"log"

	"github.com/chgrape/vaultpp/internal/repository"
	"github.com/chgrape/vaultpp/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

type UserValidator struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Role     string `json:"role" validate:"oneof=scribe member"`
}

func (s *UserService) RegisterUser(user UserValidator, ctx context.Context) (int, error) {
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
