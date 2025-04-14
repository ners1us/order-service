package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/enum"
	"github.com/ners1us/order-service/internal/metric"
	"github.com/ners1us/order-service/internal/model"
	"github.com/ners1us/order-service/internal/service"
	"net/http"
	"strconv"
	"time"
)

type PVZHandler interface {
	CreatePVZ(c *gin.Context)
	GetPVZList(c *gin.Context)
}

type pvzHandlerImpl struct {
	pvzService service.PVZService
}

func NewPVZHandler(pvzService service.PVZService) PVZHandler {
	return &pvzHandlerImpl{pvzService}
}

func (ph *pvzHandlerImpl) CreatePVZ(c *gin.Context) {
	role, _ := c.Get("role")
	var pvz model.PVZ
	if err := c.BindJSON(&pvz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdPVZ, err := ph.pvzService.CreatePVZ(&pvz, role.(string))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, enum.ErrNoModeratorRights) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	metric.PVZCreated.Inc()

	c.JSON(http.StatusCreated, createdPVZ)
}

func (ph *pvzHandlerImpl) GetPVZList(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrInvalidStartDate.Error()})
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": enum.ErrInvalidEndDate.Error()})
			return
		}
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 30 {
		limit = 10
	}

	pvzList, err := ph.pvzService.GetPVZList(startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pvzList)
}
