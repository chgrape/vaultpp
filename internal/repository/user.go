package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int
	Name      string
	Role      string
	PassHash  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) CreateUser(user User, ctx context.Context) (int, error) {
	var id int

	err := r.DB.QueryRow(ctx, "INSERT INTO users(name, password, role) VALUES($1,$2,$3) RETURNING id", user.Name, user.PassHash, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUserPassHash(name string, ctx context.Context) (*User, error) {
	var user User

	err := r.DB.QueryRow(ctx, "SELECT id, password, role FROM users WHERE users.name=$1", name).Scan(&user.ID, &user.PassHash, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
