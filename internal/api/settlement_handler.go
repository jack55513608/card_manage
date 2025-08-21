package api

import (
	"card_manage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettlementHandler struct {
	service *service.SettlementService
}

func NewSettlementHandler(settlementService *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: settlementService}
}

type CreateSettlementRequest struct {
	StoreID int64 `json:"store_id" binding:"required"`
}

// @Summary Create a new settlement request
// @Description Player creates a settlement request to clear their earnings from a store.
// @Tags settlements
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   settlement body CreateSettlementRequest true "Settlement Request Information"
// @Success 201 {object} model.Settlement
// @Failure 400 {object} map[string]string "{"error": "bad request"}"
// @Failure 409 {object} map[string]string "{"error": "conflict (e.g., no unsettled transactions)"}"
// @Failure 500 {object} map[string]string "{"error": "failed to create settlement request"}"
// @Router /api/settlements [post]
func (h *SettlementHandler) CreateSettlement(c *gin.Context) {
	var req CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	settlement, err := h.service.CreateSettlement(claims.UserID, req.StoreID)
	if err != nil {
		if err == service.ErrNoUnsettledTransactions {
			c.JSON(http.StatusConflict, gin.H{"error": "no transactions available for settlement"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create settlement request"})
		return
	}

	c.JSON(http.StatusCreated, settlement)
}

// Placeholder for store-side actions
func (h *SettlementHandler) ListSettlements(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "listing settlements (not implemented yet)"})
}

func (h *SettlementHandler) CompleteSettlement(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "completing settlement (not implemented yet)"})
}