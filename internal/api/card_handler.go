package api

import (
	"card_manage/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

type CreateCardRequest struct {
	Name       string `json:"name" binding:"required"`
	Series     string `json:"series"`
	Rarity     string `json:"rarity"`
	CardNumber string `json:"card_number"`
}

// @Summary Create a new card
// @Description Adds a new card to the store associated with the user.
// @Tags cards
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   card body CreateCardRequest true "Card Information"
// @Success 201 {object} model.Card
// @Failure 400 {object} map[string]string "{"error": "bad_request_error"}"
// @Failure 403 {object} map[string]string "{"error": "user does not have a store"}"
// @Failure 500 {object} map[string]string "{"error": "failed to create card"}"
// @Router /api/cards [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	card, err := h.cardService.CreateCard(claims.UserID, req.Name, req.Series, req.Rarity, req.CardNumber)
	if err != nil {
		if err == service.ErrStoreNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "user does not have a store"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, card)
}

// @Summary Get a card by ID
// @Description Retrieves a specific card by its ID.
// @Tags cards
// @Produce  json
// @Security BearerAuth
// @Param   id path int true "Card ID"
// @Success 200 {object} model.Card
// @Failure 400 {object} map[string]string "{"error": "invalid card ID"}"
// @Failure 403 {object} map[string]string "{"error": "you do not have permission to view this card"}"
// @Failure 404 {object} map[string]string "{"error": "card not found"}"
// @Failure 500 {object} map[string]string "{"error": "failed to retrieve card"}"
// @Router /api/cards/{id} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	cardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	card, err := h.cardService.GetCard(claims.UserID, cardID)
	if err != nil {
		if err == service.ErrCardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to view this card"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve card"})
		return
	}

	c.JSON(http.StatusOK, card)
}

// @Summary List all cards for the current user
// @Description Retrieves a list of all cards associated with the current user's store.
// @Tags cards
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} model.Card
// @Failure 500 {object} map[string]string "{"error": "failed to list cards"}"
// @Router /api/cards [get]
func (h *CardHandler) ListCards(c *gin.Context) {
	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	cards, err := h.cardService.ListCardsByCurrentUser(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

type UpdateCardRequest struct {
	Name       string `json:"name" binding:"required"`
	Series     string `json:"series"`
	Rarity     string `json:"rarity"`
	CardNumber string `json:"card_number"`
}

// @Summary Update a card
// @Description Updates the details of a specific card.
// @Tags cards
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id path int true "Card ID"
// @Param   card body UpdateCardRequest true "Card Update Information"
// @Success 200 {object} model.Card
// @Failure 400 {object} map[string]string "{"error": "invalid card ID"}"
// @Failure 403 {object} map[string]string "{"error": "you do not have permission to update this card"}"
// @Failure 404 {object} map[string]string "{"error": "card not found"}"
// @Failure 500 {object} map[string]string "{"error": "failed to update card"}"
// @Router /api/cards/{id} [put]
func (h *CardHandler) UpdateCard(c *gin.Context) {
	cardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	card, err := h.cardService.UpdateCard(claims.UserID, cardID, req.Name, req.Series, req.Rarity, req.CardNumber)
	if err != nil {
		if err == service.ErrCardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to update this card"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update card"})
		return
	}

	c.JSON(http.StatusOK, card)
}

// @Summary Delete a card
// @Description Deletes a specific card by its ID.
// @Tags cards
// @Produce  json
// @Security BearerAuth
// @Param   id path int true "Card ID"
// @Success 200 {object} map[string]string "{"message": "card deleted successfully"}"
// @Failure 400 {object} map[string]string "{"error": "invalid card ID"}"
// @Failure 403 {object} map[string]string "{"error": "you do not have permission to delete this card"}"
// @Failure 404 {object} map[string]string "{"error": "card not found"}"
// @Failure 500 {object} map[string]string "{"error": "failed to delete card"}"
// @Router /api/cards/{id} [delete]
func (h *CardHandler) DeleteCard(c *gin.Context) {
	cardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	claims := c.MustGet(AuthorizationPayloadKey).(*service.CustomClaims)

	err = h.cardService.DeleteCard(claims.UserID, cardID)
	if err != nil {
		if err == service.ErrCardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to delete this card"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card deleted successfully"})
}
