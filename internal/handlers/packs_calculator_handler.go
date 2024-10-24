package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/klemis/packs-calculator/internal/services"
	"github.com/klemis/packs-calculator/models"
	"net/http"
	"strconv"
)

type Handler struct {
	service services.PacksCalculator
}

// NewHandler creates a new Handler with the provided PacksCalculatorService.
func NewHandler(packsCalculatorService services.PacksCalculator) *Handler {
	return &Handler{
		service: packsCalculatorService,
	}
}

// AddPackSize handles adding a new pack size.
func (h *Handler) AddPackSize(c *gin.Context) {
	var req *models.PackSizeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.service.AddPackSize(req.Size); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add pack size"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pack size successfully added"})
}

// DeletePackSize handles deleting an existing pack size.
func (h *Handler) DeletePackSize(c *gin.Context) {
	var req *models.PackSizeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err := h.service.DeletePackSize(req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete pack size"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pack size successfully deleted"})
}

// CalculatePacks handles calculating the minimum packs needed for a given item quantity.
func (h *Handler) CalculatePacks(c *gin.Context) {
	quantity := c.Query("quantity")
	if quantity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No quantity parameter"})
		return
	}

	q, err := strconv.ParseUint(quantity, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity parameter"})
		return
	}

	result, err := h.service.CalculatePacks(uint32(q))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate packs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quantity": quantity,
		"packs":    result,
	})
}
