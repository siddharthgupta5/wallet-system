package services

import (
	"errors"
	"github.com/siddharthgupta5/wallet-api/internal/models"
	"github.com/siddharthgupta5/wallet-api/internal/repositories"
)

type UserService interface {
	CreateUser(name, email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(name, email string) (*models.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	user := &models.User{
		Name:  name,
		Email: email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}