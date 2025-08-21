package api

import (
	"card_manage/internal/model"
	"card_manage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

type CreateTransactionRequest struct {
	ConsignmentID int64                 `json:"consignment_id" binding:"required"`
	Price         float64               `json:"price" binding:"required,gt=0"`
	PaymentMethod model.PaymentMethod `json:"payment_method" binding:"required,oneof=CASH CREDIT"`
}

// @Summary Create a new transaction
// @Description Store creates a transaction for a sold consignment.
// @Tags transactions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   transaction body CreateTransactionRequest true "Transaction Information"
// @Success 201 {object} model.Transaction
// @Failure 400 {object} map[string]string "{"error": "bad request"}"
// @Failure 409 {object} map[string]string "{"error": "conflict (e.g., consignment not found, already sold, or forbidden)"}"
// @Failure 500 {object} map[string]string "{"error": "failed to create transaction"}"
// @Router /api/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	tx, err := h.transactionService.CreateTransaction(claims.UserID, req.ConsignmentID, req.Price, req.PaymentMethod)
	if err != nil {
		// Handle specific, known errors
		if err == service.ErrConsignmentNotFound || err == service.ErrConsignmentAlreadySold || err == service.ErrForbidden {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// Handle generic errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}