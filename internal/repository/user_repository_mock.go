package repository

import (
	"github.com/ners1us/order-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (mur *MockUserRepository) CreateUser(user *model.User) error {
	args := mur.Called(user)
	return args.Error(0)
}

func (mur *MockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	args := mur.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}
