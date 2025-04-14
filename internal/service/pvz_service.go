package service

import (
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repositories"
	"time"
)

type PVZService interface {
	CreatePVZ(pvz *models.PVZ, userRole string) (*models.PVZ, error)
	GetPVZList(startDate, endDate time.Time, page, limit int) ([]models.PVZWithReceptions, error)
}

type pvzServiceImpl struct {
	pvzRepo       repositories.PVZRepository
	receptionRepo repositories.ReceptionRepository
	productRepo   repositories.ProductRepository
}

func NewPVZService(pvzRepo repositories.PVZRepository, receptionRepo repositories.ReceptionRepository, productRepo repositories.ProductRepository) PVZService {
	return &pvzServiceImpl{
		pvzRepo,
		receptionRepo,
		productRepo,
	}
}

func (ps *pvzServiceImpl) CreatePVZ(pvz *models.PVZ, userRole string) (*models.PVZ, error) {
	if userRole != "moderator" {
		return &models.PVZ{}, enums.ErrNoModeratorRights
	}
	if pvz.City != "Москва" && pvz.City != "Санкт-Петербург" && pvz.City != "Казань" {
		return &models.PVZ{}, enums.ErrInvalidCity
	}
	if err := ps.pvzRepo.CreatePVZ(pvz); err != nil {
		return &models.PVZ{}, err
	}
	return pvz, nil
}

func (ps *pvzServiceImpl) GetPVZList(startDate, endDate time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	pvzs, err := ps.pvzRepo.GetPVZs(page, limit)
	if err != nil {
		return nil, err
	}
	pvzIDs := make([]string, len(pvzs))
	for i, pvz := range pvzs {
		pvzIDs[i] = pvz.ID
	}
	receptions, err := ps.receptionRepo.GetReceptionsByPVZIDsAndDate(pvzIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}
	receptionIDs := make([]string, len(receptions))
	for i, reception := range receptions {
		receptionIDs[i] = reception.ID
	}
	products, err := ps.productRepo.GetProductsByReceptionIDs(receptionIDs)
	if err != nil {
		return nil, err
	}

	receptionProducts := make(map[string][]models.Product)
	for _, product := range products {
		receptionProducts[product.ReceptionID] = append(receptionProducts[product.ReceptionID], product)
	}
	pvzReceptions := make(map[string][]models.ReceptionWithProducts)
	for _, reception := range receptions {
		rwp := models.ReceptionWithProducts{
			Reception: reception,
			Products:  receptionProducts[reception.ID],
		}
		pvzReceptions[reception.PVZID] = append(pvzReceptions[reception.PVZID], rwp)
	}

	var result []models.PVZWithReceptions
	for _, pvz := range pvzs {
		result = append(result, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: pvzReceptions[pvz.ID],
		})
	}
	return result, nil
}
