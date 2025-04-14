package repository

import (
	"github.com/ners1us/order-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPVZRepository struct {
	mock.Mock
}

func (mpr *MockPVZRepository) CreatePVZ(pvz *model.PVZ) error {
	args := mpr.Called(pvz)
	return args.Error(0)
}

func (mpr *MockPVZRepository) GetPVZs(page, limit int) ([]model.PVZ, error) {
	args := mpr.Called(page, limit)
	return args.Get(0).([]model.PVZ), args.Error(1)
}

func (mpr *MockPVZRepository) GetAllPVZs() ([]model.PVZ, error) {
	args := mpr.Called()
	return args.Get(0).([]model.PVZ), args.Error(1)
}

func (mpr *MockPVZRepository) GetPVZByID(id string) (*model.PVZ, error) {
	args := mpr.Called(id)
	return args.Get(0).(*model.PVZ), args.Error(1)
}
