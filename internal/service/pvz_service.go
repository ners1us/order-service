package service

import (
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"time"
)

type PVZService interface {
	CreatePVZ(pvz *model.PVZ, userRole string) (*model.PVZ, error)
	GetPVZList(startDate, endDate time.Time, page, limit int) ([]model.PVZWithReceptions, error)
}

type pvzServiceImpl struct {
	pvzRepo       repository.PVZRepository
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
}

func NewPVZService(pvzRepo repository.PVZRepository, receptionRepo repository.ReceptionRepository, productRepo repository.ProductRepository) PVZService {
	return &pvzServiceImpl{
		pvzRepo,
		receptionRepo,
		productRepo,
	}
}

func (ps *pvzServiceImpl) CreatePVZ(pvz *model.PVZ, userRole string) (*model.PVZ, error) {
	if userRole != "moderator" {
		return &model.PVZ{}, enum.ErrNoModeratorRights
	}
	if pvz.City != "Москва" && pvz.City != "Санкт-Петербург" && pvz.City != "Казань" {
		return &model.PVZ{}, enum.ErrInvalidCity
	}
	if err := ps.pvzRepo.CreatePVZ(pvz); err != nil {
		return &model.PVZ{}, err
	}
	return pvz, nil
}

func (ps *pvzServiceImpl) GetPVZList(startDate, endDate time.Time, page, limit int) ([]model.PVZWithReceptions, error) {
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

	receptionProducts := make(map[string][]model.Product)
	for _, product := range products {
		receptionProducts[product.ReceptionID] = append(receptionProducts[product.ReceptionID], product)
	}
	pvzReceptions := make(map[string][]model.ReceptionWithProducts)
	for _, reception := range receptions {
		rwp := model.ReceptionWithProducts{
			Reception: reception,
			Products:  receptionProducts[reception.ID],
		}
		pvzReceptions[reception.PVZID] = append(pvzReceptions[reception.PVZID], rwp)
	}

	var result []model.PVZWithReceptions
	for _, pvz := range pvzs {
		result = append(result, model.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: pvzReceptions[pvz.ID],
		})
	}
	return result, nil
}
