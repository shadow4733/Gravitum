package service

import (
	"Gravitum/internal/model"
	"Gravitum/internal/repo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(firstName, lastName, email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
