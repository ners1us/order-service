package repository

import (
	"github.com/ners1us/order-service/internal/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockReceptionRepository struct {
	mock.Mock
}

func (mrr *MockReceptionRepository) CreateReception(reception *models.Reception) error {
	args := mrr.Called(reception)
	return args.Error(0)
}

func (mrr *MockReceptionRepository) GetLastReceptionByPVZID(pvzID string) (*models.Reception, error) {
	args := mrr.Called(pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (mrr *MockReceptionRepository) UpdateReceptionStatus(id string, status string) error {
	args := mrr.Called(id, status)
	return args.Error(0)
}

func (mrr *MockReceptionRepository) GetReceptionsByPVZIDsAndDate(pvzIDs []string, startDate, endDate time.Time) ([]models.Reception, error) {
	args := mrr.Called(pvzIDs, startDate, endDate)
	return args.Get(0).([]models.Reception), args.Error(1)
}
