package services

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repositories"
	"time"
)

type ReceptionService interface {
	CreateReception(pvzID string, userRole string) (*models.Reception, error)
	CloseLastReception(pvzID string, userRole string) (*models.Reception, error)
}

type receptionServiceImpl struct {
	receptionRepo repositories.ReceptionRepository
	pvzRepo       repositories.PVZRepository
}

func NewReceptionService(receptionRepo repositories.ReceptionRepository, pvzRepo repositories.PVZRepository) ReceptionService {
	return &receptionServiceImpl{
		receptionRepo,
		pvzRepo,
	}
}

func (rs *receptionServiceImpl) CreateReception(pvzID string, userRole string) (*models.Reception, error) {
	if userRole != "employee" {
		return &models.Reception{}, enums.ErrNoEmployeeRights
	}

	pvz, err := rs.pvzRepo.GetPVZByID(pvzID)
	if err != nil {
		return &models.Reception{}, err
	}
	if pvz.ID == "" {
		return &models.Reception{}, enums.ErrPVZNotFound
	}

	lastReception, err := rs.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &models.Reception{}, err
	}
	if lastReception.Status == "in_progress" {
		return &models.Reception{}, enums.ErrOpenReception
	}

	reception := models.Reception{
		ID:       uuid.New().String(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}
	if err := rs.receptionRepo.CreateReception(&reception); err != nil {
		return &models.Reception{}, err
	}
	return &reception, nil
}

func (rs *receptionServiceImpl) CloseLastReception(pvzID string, userRole string) (*models.Reception, error) {
	if userRole != "employee" {
		return &models.Reception{}, enums.ErrNoEmployeeRights
	}
	lastReception, err := rs.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &models.Reception{}, err
	}
	if lastReception.Status != "in_progress" {
		return &models.Reception{}, enums.ErrNoOpenReceptionToClose
	}
	if err := rs.receptionRepo.UpdateReceptionStatus(lastReception.ID, "closed"); err != nil {
		return &models.Reception{}, err
	}
	lastReception.Status = "closed"
	return lastReception, nil
}
