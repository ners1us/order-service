package service

import (
	"errors"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatePVZ_InvalidRole(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := new(model.PVZ)
	userRole := enum.RoleEmployee.String()

	// Act
	_, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoModeratorRights, err)
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := &model.PVZ{City: "InvalidCity"}
	userRole := enum.RoleModerator.String()

	// Act
	_, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrInvalidCity, err)
}

func TestCreatePVZ_ValidCitySPb(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)
	pvz := &model.PVZ{City: enum.CitySaintPetersburg.String()}
	userRole := enum.RoleModerator.String()
	mockPVZRepo.On("CreatePVZ", pvz).Return(nil)

	// Act
	result, err := service.CreatePVZ(pvz, userRole)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, pvz, result)
}

func TestGetPVZList_Success(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	page, limit := 1, 10
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	pvzs := []model.PVZ{
		{ID: "pvz_1", City: enum.CityMoscow.String()},
		{ID: "pvz_2", City: enum.CitySaintPetersburg.String()},
	}
	pvzIDs := []string{"pvz_1", "pvz_2"}

	receptions := []model.Reception{
		{ID: "rec_1", PVZID: "pvz_1", Status: enum.StatusClosed.String()},
		{ID: "rec_2", PVZID: "pvz_2", Status: enum.StatusInProgress.String()},
	}
	receptionIDs := []string{"rec_1", "rec_2"}

	products := []model.Product{
		{ID: "prod_1", ReceptionID: "rec_1", Type: "type_1"},
		{ID: "prod_2", ReceptionID: "rec_2", Type: "type_2"},
	}

	mockPVZRepo.On("GetPVZs", page, limit).Return(pvzs, nil)
	mockReceptionRepo.On("GetReceptionsByPVZIDsAndDate", pvzIDs, startDate, endDate).Return(receptions, nil)
	mockProductRepo.On("GetProductsByReceptionIDs", receptionIDs).Return(products, nil)

	// Act
	result, err := service.GetPVZList(startDate, endDate, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "pvz_1", result[0].PVZ.ID)
	assert.Equal(t, "pvz_2", result[1].PVZ.ID)
	assert.Len(t, result[0].Receptions, 1)
	assert.Len(t, result[1].Receptions, 1)
	assert.Equal(t, "rec_1", result[0].Receptions[0].Reception.ID)
	assert.Equal(t, "rec_2", result[1].Receptions[0].Reception.ID)
	assert.Len(t, result[0].Receptions[0].Products, 1)
	assert.Len(t, result[1].Receptions[0].Products, 1)
	assert.Equal(t, "prod_1", result[0].Receptions[0].Products[0].ID)
	assert.Equal(t, "prod_2", result[1].Receptions[0].Products[0].ID)
}

func TestGetPVZList_PVZRepoError(t *testing.T) {
	// Arrange
	mockPVZRepo := new(repository.MockPVZRepository)
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	page, limit := 1, 10
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	mockPVZRepo.On("GetPVZs", page, limit).Return([]model.PVZ{}, errors.New("db error"))

	// Act
	result, err := service.GetPVZList(startDate, endDate, page, limit)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "db error", err.Error())
}
