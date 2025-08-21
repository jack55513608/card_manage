package api

import (
	"card_manage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler(storeService *service.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

type CreateStoreRequest struct {
	Name             string  `json:"name" binding:"required"`
	CommissionCash   float64 `json:"commission_cash" binding:"gte=0,lte=100"`
	CommissionCredit float64 `json:"commission_credit" binding:"gte=0,lte=100"`
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user claims from the context (set by AuthMiddleware)
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	claims := payload.(*service.CustomClaims)

	store, err := h.storeService.CreateStore(claims.UserID, req.Name, req.CommissionCash, req.CommissionCredit)
	if err != nil {
		// In a real app, you'd want more sophisticated error handling
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create store"})
		return
	}

	c.JSON(http.StatusCreated, store)
}
