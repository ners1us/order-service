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

func TestCreateReception_NoEmployee(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := enum.RoleModerator.String()

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrNoEmployeeRights, err)
}

func TestCloseLastReception_Success(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()
	lastReception := &model.Reception{ID: "rec_1", Status: enum.StatusInProgress.String()}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockReceptionRepo.On("UpdateReceptionStatus", "rec_1", enum.StatusClosed.String()).Return(nil)

	// Act
	result, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, enum.StatusClosed.String(), result.Status)
}

func TestCreateReception_PVZNotFound(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()
	mockPVZRepo.On("GetPVZByID", pvzID).Return(&model.PVZ{}, nil)

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enum.ErrPVZNotFound, err)
}

func TestCreateReception_RepoError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()
	pvz := &model.PVZ{ID: "test_pvz_id"}
	mockPVZRepo.On("GetPVZByID", pvzID).Return(pvz, nil)
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&model.Reception{Status: enum.StatusClosed.String()}, nil)
	mockReceptionRepo.On("CreateReception", mock.Anything).Return(errors.New("create error"))

	// Act
	_, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "create error", err.Error())
}

func TestCloseLastReception_UpdateError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)
	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()
	lastReception := &model.Reception{ID: "rec_1", Status: enum.StatusInProgress.String()}
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(lastReception, nil)
	mockReceptionRepo.On("UpdateReceptionStatus", "rec_1", enum.StatusClosed.String()).Return(errors.New("update error"))

	// Act
	_, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
}

func TestCreateReception_GetLastReceptionError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)

	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()

	mockPVZRepo.On("GetPVZByID", pvzID).Return(&model.PVZ{ID: pvzID}, nil)
	mockReceptionRepo.On("GetLastReceptionByPVZID", pvzID).Return(&model.Reception{}, errors.New("reception error"))

	// Act
	result, err := service.CloseLastReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "reception error", err.Error())
	assert.Empty(t, result.ID)
}

func TestCreateReception_GetPVZError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(repository.MockReceptionRepository)
	mockPVZRepo := new(repository.MockPVZRepository)
	service := NewReceptionService(mockReceptionRepo, mockPVZRepo)

	pvzID := "test_pvz_id"
	userRole := enum.RoleEmployee.String()

	mockPVZRepo.On("GetPVZByID", pvzID).Return(&model.PVZ{}, errors.New("PVZ error"))

	// Act
	result, err := service.CreateReception(pvzID, userRole)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "PVZ error", err.Error())
	assert.Empty(t, result.ID)
}
