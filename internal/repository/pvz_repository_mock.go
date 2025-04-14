package repository

import (
	"github.com/ners1us/order-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockPVZRepository struct {
	mock.Mock
}

func (mpr *MockPVZRepository) CreatePVZ(pvz *models.PVZ) error {
	args := mpr.Called(pvz)
	return args.Error(0)
}

func (mpr *MockPVZRepository) GetPVZs(page, limit int) ([]models.PVZ, error) {
	args := mpr.Called(page, limit)
	return args.Get(0).([]models.PVZ), args.Error(1)
}

func (mpr *MockPVZRepository) GetAllPVZs() ([]models.PVZ, error) {
	args := mpr.Called()
	return args.Get(0).([]models.PVZ), args.Error(1)
}

func (mpr *MockPVZRepository) GetPVZByID(id string) (*models.PVZ, error) {
	args := mpr.Called(id)
	return args.Get(0).(*models.PVZ), args.Error(1)
}
