package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/metric"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/service"
	"net/http"
)

type ProductHandler interface {
	AddProduct(c *gin.Context)
	DeleteLastProduct(c *gin.Context)
}

type productHandlerImpl struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandlerImpl{productService}
}

func (ph *productHandlerImpl) AddProduct(c *gin.Context) {
	role, _ := c.Get("role")
	var req struct {
		Type  string `json:"type"`
		PVZID string `json:"pvzId"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := model.Product{Type: req.Type}
	createdProduct, err := ph.productService.AddProduct(&product, req.PVZID, role.(string))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, enum.ErrNoEmployeeRights) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	metric.ProductsAdded.Inc()

	c.JSON(http.StatusCreated, createdProduct)
}

func (ph *productHandlerImpl) DeleteLastProduct(c *gin.Context) {
	pvzID := c.Param("pvzId")
	role, _ := c.Get("role")
	if err := ph.productService.DeleteLastProduct(pvzID, role.(string)); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, enum.ErrNoEmployeeRights) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
