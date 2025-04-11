package services

import (
	"github.com/google/uuid"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/models"
	"github.com/ners1us/order-service/internal/repositories"
	"time"
)

type ProductService interface {
	AddProduct(product *models.Product, pvzID string, userRole string) (*models.Product, error)
	DeleteLastProduct(pvzID string, userRole string) error
}

type productServiceImpl struct {
	receptionRepo repositories.ReceptionRepository
	productRepo   repositories.ProductRepository
}

func NewProductService(receptionRepo repositories.ReceptionRepository, productRepo repositories.ProductRepository) ProductService {
	return &productServiceImpl{
		receptionRepo,
		productRepo,
	}
}

func (ps *productServiceImpl) AddProduct(product *models.Product, pvzID string, userRole string) (*models.Product, error) {
	if userRole != "employee" {
		return &models.Product{}, enums.ErrNoEmployeeRights
	}
	lastReception, err := ps.receptionRepo.GetLastReceptionByPVZID(pvzID)
	if err != nil {
		return &models.Product{}, err
	}
	if lastReception.Status != "in_progress" {
		return &models.Product{}, enums.ErrNoOpenReceptionsToAdd
	}
	product.ID = uuid.New().String()
	product.DateTime = time.Now()
	product.ReceptionID = lastReception.ID
	if err := ps.productRepo.CreateProduct(product); err != nil {
		return &models.Product{}, err
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
