package api

import (
	"card_manage/internal/model"
	"card_manage/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConsignmentHandler struct {
	consignmentService *service.ConsignmentService
}

func NewConsignmentHandler(consignmentService *service.ConsignmentService) *ConsignmentHandler {
	return &ConsignmentHandler{consignmentService: consignmentService}
}

type CreateConsignmentRequest struct {
	StoreID int64   `json:"store_id" binding:"required"`
	CardIDs []int64 `json:"card_ids" binding:"required,gt=0"`
}

// @Summary Create a new consignment request
// @Description Player creates a consignment request for one or more cards to a store.
// @Tags consignments
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   consignment body CreateConsignmentRequest true "Consignment Request Information"
// @Success 201 {object} model.Consignment
// @Failure 400 {object} map[string]string "{"error": "bad request"}"
// @Failure 500 {object} map[string]string "{"error": "failed to create consignment request"}"
// @Router /api/consignments [post]
func (h *ConsignmentHandler) CreateConsignment(c *gin.Context) {
	var req CreateConsignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	consignment, err := h.consignmentService.CreateConsignment(claims.UserID, req.StoreID, req.CardIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create consignment request"})
		return
	}

	c.JSON(http.StatusCreated, consignment)
}

type UpdateConsignmentItemStatusRequest struct {
	Status model.ConsignmentItemStatus `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	Reason string                      `json:"reason"`
}

// @Summary Update a consignment item's status
// @Description Store approves or rejects a consignment item.
// @Tags consignments
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   itemId path int true "Consignment Item ID"
// @Param   status body UpdateConsignmentItemStatusRequest true "New Status"
// @Success 200 {object} model.ConsignmentItem
// @Failure 400 {object} map[string]string "{"error": "invalid item ID or bad request"}"
// @Failure 403 {object} map[string]string "{"error": "permission denied"}"
// @Failure 404 {object} map[string]string "{"error": "item not found"}"
// @Failure 500 {object} map[string]string "{"error": "failed to update item status"}"
// @Router /api/consignments/items/{itemId} [put]
func (h *ConsignmentHandler) UpdateConsignmentItemStatus(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	var req UpdateConsignmentItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	item, err := h.consignmentService.UpdateConsignmentItemStatus(claims.UserID, itemID, req.Status, req.Reason)
	if err != nil {
		switch err {
		case service.ErrConsignmentItemNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		case service.ErrCannotUpdateStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item status"})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}
