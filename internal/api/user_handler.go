package api

import (
	"card_manage/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	jwtService  *service.JWTService
}

func NewUserHandler(userService *service.UserService, jwtService *service.JWTService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=PLAYER STORE"` // ADMIN role cannot be self-assigned
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary User Registration
// @Description Creates a new user account.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body RegisterRequest true "User Registration Info"
// @Success 201 {object} map[string]interface{} "{"message": "user created successfully", "user_id": 1}"
// @Failure 400 {object} map[string]string "{"error": "bad_request_error"}"
// @Failure 409 {object} map[string]string "{"error": "email already exists"}"
// @Failure 500 {object} map[string]string "{"error": "failed to register user"}"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		case errors.Is(err, service.ErrDatabase):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user_id": user.ID})
}

// @Summary User Login
// @Description Authenticates a user and returns a JWT token.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   credentials body LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]string "{"token": "your_jwt_token"}"
// @Failure 400 {object} map[string]string "{"error": "bad_request_error"}"
// @Failure 401 {object} map[string]string "{"error": "invalid email or password"}"
// @Failure 500 {object} map[string]string "{"error": "internal_server_error"}"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		case errors.Is(err, service.ErrDatabase):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get user by ID
// @Description Get user details by ID. Requires authentication and ADMIN role.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "{"error": "invalid user ID"}"
// @Failure 401 {object} map[string]string "{"error": "unauthorized"}"
// @Failure 403 {object} map[string]string "{"error": "forbidden"}"
// @Failure 404 {object} map[string]string "{"error": "user not found"}"
// @Failure 500 {object} map[string]string "{"error": "database error"}"
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, service.ErrDatabase):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}