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
	StoreID  int64 `json:"store_id" binding:"required"`
	CardID   int64 `json:"card_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,gt=0"`
}

// @Summary Create a new consignment request
// @Description Player creates a consignment request for a card to a store.
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

	consignment, err := h.consignmentService.CreateConsignment(claims.UserID, req.StoreID, req.CardID, req.Quantity)
	if err != nil {
		if err == service.ErrInvalidCardForStore {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the provided card does not belong to the specified store"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create consignment request"})
		return
	}

	c.JSON(http.StatusCreated, consignment)
}

// @Summary List consignments
// @Description Lists consignments for the current user (player or store).
// @Tags consignments
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} model.Consignment
// @Failure 500 {object} map[string]string "{"error": "failed to list consignments"}"
// @Router /api/consignments [get]
func (h *ConsignmentHandler) ListConsignments(c *gin.Context) {
	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	consignments, err := h.consignmentService.ListConsignmentsForUser(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list consignments"})
		return
	}

	c.JSON(http.StatusOK, consignments)
}

type UpdateConsignmentStatusRequest struct {
	Status model.ConsignmentStatus `json:"status" binding:"required,oneof=LISTED SOLD CLEARED"`
}

// @Summary Update a consignment's status
// @Description Store updates the status of a consignment (e.g., to LISTED, SOLD).
// @Tags consignments
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id path int true "Consignment ID"
// @Param   status body UpdateConsignmentStatusRequest true "New Status"
// @Success 200 {object} model.Consignment
// @Failure 400 {object} map[string]string "{"error": "invalid consignment ID or bad request"}"
// @Failure 403 {object} map[string]string "{"error": "permission denied"}"
// @Failure 404 {object} map[string]string "{"error": "consignment not found"}"
// @Failure 500 {object} map[string]string "{"error": "failed to update consignment status"}"
// @Router /api/consignments/{id} [put]
func (h *ConsignmentHandler) UpdateConsignmentStatus(c *gin.Context) {
	consignmentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid consignment ID"})
		return
	}

	var req UpdateConsignmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	consignment, err := h.consignmentService.UpdateConsignmentStatus(claims.UserID, consignmentID, req.Status)
	if err != nil {
		if err == service.ErrConsignmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "consignment not found"})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to update this consignment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update consignment status"})
		return
	}

	c.JSON(http.StatusOK, consignment)
}