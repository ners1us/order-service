package service

import (
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePVZ_InvalidRole(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := new(models.PVZ)
	userRole := "employee"

	// Act
	_, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enums.ErrNoModeratorRights, err)
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := &models.PVZ{City: "InvalidCity"}
	userRole := "moderator"

	// Act
	_, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enums.ErrInvalidCity, err)
}

func TestCreatePVZ_ValidCitySPb(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := &models.PVZ{City: "Санкт-Петербург"}
	userRole := "moderator"
	mockPVZRepo.On("CreatePVZ", pvz).Return(nil)

	// Act
	result, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pvz, result)
}
