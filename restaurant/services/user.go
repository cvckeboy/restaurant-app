package services

import (
	"context"
	"errors"
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/restaurant/storage"
	"github.com/cvckeboy/restaurant-app/utils"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) RegisterUser(ctx context.Context, req *models.RegisterUserRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     req.Role,
	}
	return s.storage.RegisterUser(ctx, user)
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *models.LoginUserRequest) (*models.User, error) {
	user, err := s.storage.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
