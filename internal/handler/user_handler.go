package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/auth" // Ganti dengan path modul Anda
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

// Register menangani permintaan registrasi pengguna baru
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	// 1. Bind and Validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Create user model for database
	newUser := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// 4. Save to database
	userID, err := h.userStore.Create(c.Request.Context(), newUser)
	if err != nil {
		// Cek error duplikat (sangat umum)
		// Kode error '23505' adalah unique_violation di PostgreSQL
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

    // Untuk mendapatkan data lengkap (termasuk created_at), idealnya kita SELECT lagi
    // Tapi untuk sekarang, kita buat response manual saja
	response := model.UserResponse{
		ID:       userID,
		Username: newUser.Username,
		Email:    newUser.Email,
        // CreatedAt akan diisi oleh database, ini hanya contoh
	}

	// 5. Send success response
	c.JSON(http.StatusCreated, response)
}
