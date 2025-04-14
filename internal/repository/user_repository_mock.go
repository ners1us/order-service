package repository

import (
	"github.com/ners1us/order-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (mur *MockUserRepository) CreateUser(user *models.User) error {
	args := mur.Called(user)
	return args.Error(0)
}

func (mur *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := mur.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}
