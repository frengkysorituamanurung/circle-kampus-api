package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/auth" 
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/model"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/store"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	userStore *store.UserStore
	validate  *validator.Validate
}

func NewUserHandler(userStore *store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		validate:  validator.New(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userID, err := h.userStore.Create(c.Request.Context(), newUser)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := model.UserResponse{
		ID:       userID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userStore.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan internal"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{AccessToken: token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Ambil userID dari context yang sudah di-set oleh middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID from context"})
		return
	}

	// Panggil store untuk mendapatkan detail user dari DB
	user, err := h.userStore.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user profile"})
		return
	}

	// Kembalikan response profil
	response := model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}