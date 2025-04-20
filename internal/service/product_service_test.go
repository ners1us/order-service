package service

import (
	"errors"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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
