package service

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"time"
)

type ReceptionService interface {
	CreateReception(pvzID string, userRole string) (*model.Reception, error)
	CloseLastReception(pvzID string, userRole string) (*model.Reception, error)
}

type receptionServiceImpl struct {
	receptionRepo repository.ReceptionRepository
	pvzRepo       repository.PVZRepository
}

func NewReceptionService(receptionRepo repository.ReceptionRepository, pvzRepo repository.PVZRepository) ReceptionService {
	return &receptionServiceImpl{
		receptionRepo,
		pvzRepo,
	}
}

func (rs *receptionServiceImpl) CreateReception(pvzID string, userRole string) (*model.Reception, error) {
	if userRole != "employee" {
		return &model.Reception{}, enums.ErrNoEmployeeRights
	}

	pvz, err := rs.pvzRepo.GetPVZByID(pvzID)
	if err != nil {
		return &model.Reception{}, err
	}
	if pvz.ID == "" {
		return &model.Reception{}, enums.ErrPVZNotFound
	}

	lastReception, err := rs.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &model.Reception{}, err
	}
	if lastReception.Status == "in_progress" {
		return &model.Reception{}, enums.ErrOpenReception
	}

	reception := model.Reception{
		ID:       uuid.New().String(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}
	if err := rs.receptionRepo.CreateReception(&reception); err != nil {
		return &model.Reception{}, err
	}
	return &reception, nil
}

func (rs *receptionServiceImpl) CloseLastReception(pvzID string, userRole string) (*model.Reception, error) {
	if userRole != "employee" {
		return &model.Reception{}, enums.ErrNoEmployeeRights
	}
	lastReception, err := rs.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &model.Reception{}, err
	}
	if lastReception.Status != "in_progress" {
		return &model.Reception{}, enums.ErrNoOpenReceptionToClose
	}
	if err := rs.receptionRepo.UpdateReceptionStatus(lastReception.ID, "closed"); err != nil {
		return &model.Reception{}, err
	}
	lastReception.Status = "closed"
	return lastReception, nil
}
