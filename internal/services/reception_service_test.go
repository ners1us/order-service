package services

import (
	"errors"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateReception_NoEmployee(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := "moderator"

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enums.ErrNoEmployeeRights, err)
}

func TestCloseLastReception_Success(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	lastReception := &models.Reception{ID: "rec_1", Status: "in_progress"}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockReceptionRepo.On("UpdateReceptionStatus", "rec_1", "closed").Return(nil)

	// Act
	result, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "closed", result.Status)
}

func TestCreateReception_PVZNotFound(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	mockPVZRepo.On("GetPVZByID", pvzID).Return(&models.PVZ{}, nil)

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enums.ErrPVZNotFound, err)
}

func TestCreateReception_RepoError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	pvz := &models.PVZ{ID: "test_pvz_id"}
	mockPVZRepo.On("GetPVZByID", pvzID).Return(pvz, nil)
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&models.Reception{Status: "closed"}, nil)
	mockReceptionRepo.On("CreateReception", mock.Anything).Return(errors.New("create error"))

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "create error", err.Error())
}

func TestCloseLastReception_UpdateError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := "employee"
	lastReception := &models.Reception{ID: "rec_1", Status: "in_progress"}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockReceptionRepo.On("UpdateReceptionStatus", "rec_1", "closed").Return(errors.New("update error"))

	// Act
	_, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
}

func TestCreateReception_GetLastReceptionError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)

	pvzID := "test_pvz_id"
	userRole := "employee"

	mockPVZRepo.On("GetPVZByID", pvzID).Return(&models.PVZ{ID: pvzID}, nil)
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&models.Reception{}, errors.New("reception error"))

	// Act
	result, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "reception error", err.Error())
	assert.Empty(t, result.ID)
}

func TestCreateReception_GetPVZError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repositories.MockReceptionRepository)
	mockPVZRepo := new(repositories.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)

	pvzID := "test_pvz_id"
	userRole := "employee"

	mockPVZRepo.On("GetPVZByID", pvzID).Return(&models.PVZ{}, errors.New("PVZ error"))

	// Act
	result, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "PVZ error", err.Error())
	assert.Empty(t, result.ID)
}
