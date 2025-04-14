package service

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/repository"
	"time"
)

type ProductService interface {
	AddProduct(product *model.Product, pvzID string, userRole string) (*model.Product, error)
	DeleteLastProduct(pvzID string, userRole string) error
}

type productServiceImpl struct {
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
}

func NewProductService(receptionRepo repository.ReceptionRepository, productRepo repository.ProductRepository) ProductService {
	return &productServiceImpl{
		receptionRepo,
		productRepo,
	}
}

func (ps *productServiceImpl) AddProduct(product *model.Product, pvzID string, userRole string) (*model.Product, error) {
	if userRole != "employee" {
		return &model.Product{}, enums.ErrNoEmployeeRights
	}
	lastReception, err := ps.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &model.Product{}, err
	}
	if lastReception.Status != "in_progress" {
		return &model.Product{}, enums.ErrNoOpenReceptionsToAdd
	}
	product.ID = uuid.New().String()
	product.DateTime = time.Now()
	product.ReceptionID = lastReception.ID
	if err := ps.productRepo.CreateProduct(product); err != nil {
		return &model.Product{}, err
	}
	return product, nil
}

func (ps *productServiceImpl) DeleteLastProduct(pvzID string, userRole string) error {
	if userRole != "employee" {
		return enums.ErrNoEmployeeRights
	}
	lastReception, err := ps.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return err
	}
	if lastReception.Status != "in_progress" {
		return enums.ErrNoOpenReceptionToDelete
	}
	lastProduct, err := ps.productRepo.GetLastProductByReceptionID(lastReception.ID)
	if err != nil {
		return err
	}
	if lastProduct.ID == "" {
		return enums.ErrNoProductsToDelete
	}
	return ps.productRepo.DeleteProduct(lastProduct.ID)
}
