package storage

import (
	"context"
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{pool: pool}
}

func (s *UserStorage) RegisterUser(ctx context.Context, user *models.User) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO users (id, username, password, role) VALUES ($1, $2, $3, $4)",
		uuid.New(), user.Username, user.Password, user.Role)
	return err
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.pool.QueryRow(ctx, "SELECT id, username, password, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
