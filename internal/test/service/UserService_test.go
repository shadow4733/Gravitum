package service

import (
	"Gravitum/internal/model"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
}

func NewUserService(repo UserRepository) *UserService {
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

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	firstName := "Ivan"
	lastName := "Ivanov"
	email := "IvanIvanov@example.com"
	password := "password123"

	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	user, err := userService.CreateUser(firstName, lastName, email, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
	assert.Equal(t, email, user.Email)
	assert.NotEmpty(t, user.ID)
	assert.True(t, time.Since(user.CreatedAt) < time.Minute)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &model.User{
		ID:        uuid.NewString(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "hashedPassword123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("*model.User")).Return(nil)

	user.FirstName = "Jonathan"

	updatedUser, err := userService.UpdateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "Jonathan", updatedUser.FirstName)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userID := uuid.NewString()

	mockRepo.On("GetUserByID", userID).Return(&model.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "hashedPassword123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	user, err := userService.GetUserByID(userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "johndoe@example.com", user.Email)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Failure(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(fmt.Errorf("some error"))

	user, err := userService.CreateUser("John", "Doe", "johndoe@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Failure(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &model.User{
		ID:        uuid.NewString(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "hashedPassword123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("*model.User")).Return(fmt.Errorf("some error"))

	updatedUser, err := userService.UpdateUser(user)

	assert.Error(t, err)
	assert.Nil(t, updatedUser)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Failure(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userID := uuid.NewString()

	mockRepo.On("GetUserByID", userID).Return(nil, fmt.Errorf("user not found"))

	user, err := userService.GetUserByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}
