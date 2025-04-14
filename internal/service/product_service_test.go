package service

import (
	"errors"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAddProduct_Success(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	product := new(model.Product)
	pvzID := "test_pvz_id"
	userRole := "employee"
	lastReception := &model.Reception{ID: "rec_1", Status: "in_progress"}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockProductRepo.On("CreateProduct", mock.Anything).Return(nil)

	// Act
	result, err := service.AddProduct(product, pvzID, userRole)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, lastReception.ID, result.ReceptionID)
}

func TestAddProduct_ProductRepoError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)

	product := &model.Product{}
	pvzID := "test_pvz_id"
	userRole := "employee"

	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).
		Return(&model.Reception{ID: "rec_1", Status: "in_progress"}, nil)
	mockProductRepo.On("CreateProduct", mock.Anything).Return(errors.New("product error"))

	// Act
	_, err := service.AddProduct(product, pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product error", err.Error())
}

func TestAddProduct_NotEmployee(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	product := new(model.Product)
	pvzID := "test_pvz_id2"
	userRole := "moderator"

	// Act
	_, err := service.AddProduct(product, pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoEmployeeRights, err)
}

func TestAddProduct_NoOpenReception(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	product := new(model.Product)
	pvzID := "test_pvz_id"
	userRole := "employee"
	lastReception := &model.Reception{Status: "closed"}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)

	// Act
	_, err := service.AddProduct(product, pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoOpenReceptionsToAdd, err)
}

func TestDeleteLastProduct_InvalidRole(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	pvzID := "test_pvz_id"
	userRole := "test_user_id"

	// Act
	err := service.DeleteLastProduct(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoEmployeeRights, err)
}

func TestDeleteLastProduct_EmptyReceptionID(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&model.Reception{ID: ""}, nil)

	// Act
	err := service.DeleteLastProduct(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoOpenReceptionToDelete, err)
}

func TestAddProduct_EmptyPVZID(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	product := new(model.Product)
	pvzID := ""
	userRole := "employee"
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&model.Reception{}, nil)

	// Act
	_, err := service.AddProduct(product, pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoOpenReceptionsToAdd, err)
}

func TestDeleteLastProduct_ProductRepoError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	lastReception := &model.Reception{ID: "rec_1", Status: "in_progress"}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockProductRepo.On("GetLastProductByReceptionID", lastReception.ID).Return(&model.Product{ID: "prod_1"}, nil)
	mockProductRepo.On("DeleteProduct", "prod_1").Return(errors.New("delete error"))

	// Act
	err := service.DeleteLastProduct(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
}

func TestDeleteLastProduct_GetLastProductError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockProductRepo := new(repository.MockProductRepository)
	service := NewProductService(mockReceptionRepo, mockProductRepo)

	pvzID := "test_pvz_id"
	userRole := "employee"
	receptionID := "rec_77"

	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).
		Return(&model.Reception{ID: receptionID, Status: "in_progress"}, nil)
	mockProductRepo.On("GetLastProductByReceptionID", receptionID).
		Return(&model.Product{}, errors.New("product error"))

	// Act
	err := service.DeleteLastProduct(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product error", err.Error())
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
		{ID: "pvz_1", City: "Москва"},
		{ID: "pvz_2", City: "Санкт-Петербург"},
	}
	pvzIDs := []string{"pvz_1", "pvz_2"}

	receptions := []model.Reception{
		{ID: "rec_1", PVZID: "pvz_1", Status: "closed"},
		{ID: "rec_2", PVZID: "pvz_2", Status: "in_progress"},
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
