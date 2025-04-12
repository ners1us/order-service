package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/enums"
	"github.com/ners1us/order-service/internal/services"
	"net/http"
)

type ReceptionHandler interface {
	CreateReception(c *gin.Context)
	CloseLastReception(c *gin.Context)
}

type receptionHandlerImpl struct {
	receptionService services.ReceptionService
}

func NewReceptionHandler(receptionService services.ReceptionService) ReceptionHandler {
	return &receptionHandlerImpl{receptionService}
}

func (rh *receptionHandlerImpl) CreateReception(c *gin.Context) {
	role, _ := c.Get("role")
	var req struct {
		PVZID string `json:"pvzId"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reception, err := rh.receptionService.CreateReception(req.PVZID, role.(string))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, enums.ErrNoModeratorRights) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reception)
}

func (rh *receptionHandlerImpl) CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")
	role, _ := c.Get("role")
	reception, err := rh.receptionService.CloseLastReception(pvzID, role.(string))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, enums.ErrNoModeratorRights) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reception)
}
