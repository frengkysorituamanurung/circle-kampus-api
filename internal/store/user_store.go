package store

import (
	"context"
	"errors"

	"github.com/frengkysorituamanurung/circle-kampus-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	//"google.golang.org/protobuf/internal/errors"
)

type UserStore struct {
	db *pgxpool.Pool
}

func (s *UserStore) Create(context context.Context, newUser *model.User) (any, error) {
	panic("unimplemented")
}

var ErrUserNotFound = errors.New("user not found")

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE email = $1`

	var user model.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
