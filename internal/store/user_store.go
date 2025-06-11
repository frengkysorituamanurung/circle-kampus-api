package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool" // Menggunakan pgxpool untuk connection pool
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/model" // Ganti dengan path modul Anda
)

type UserStore struct {
	db *pgxpool.Pool
}

// NewUserStore membuat instance UserStore baru
func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

// Create menyisipkan user baru ke dalam database
func (s *UserStore) Create(ctx context.Context, user *model.User) (string, error) {
	query := `INSERT INTO users (username, email, password_hash)
			   VALUES ($1, $2, $3)
			   RETURNING id`

	var userID string
	err := s.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		// Nanti kita bisa handle error spesifik, e.g., duplicate email/username
		return "", err
	}

	return userID, nil
}