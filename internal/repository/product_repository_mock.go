package repository

import (
	"github.com/ners1us/order-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (mpr *MockProductRepository) CreateProduct(product *model.Product) error {
	args := mpr.Called(product)
	return args.Error(0)
}

func (mpr *MockProductRepository) GetLastProductByReceptionID(receptionID string) (*model.Product, error) {
	args := mpr.Called(receptionID)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (mpr *MockProductRepository) DeleteProduct(id string) error {
	args := mpr.Called(id)
	return args.Error(0)
}

func (mpr *MockProductRepository) GetProductsByReceptionIDs(receptionIDs []string) ([]model.Product, error) {
	args := mpr.Called(receptionIDs)
	return args.Get(0).([]model.Product), args.Error(1)
}
