package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"microservices-go/services/user/internal/domain"
	"strings"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicate      = errors.New("record duplication")
	ErrEditConflict   = errors.New("edit conflict")
)

type userRepo struct {
	db *pgxpool.Pool
}

type User interface {
	Insert(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (s *userRepo) Insert(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	args := []any{user.Name, user.Email, user.HashPassword}

	err := s.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_unique"):
			return ErrDuplicate
		default:
			return err
		}
	}
	return nil
}

func (s *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, name, email, password_hash, created_at
	FROM users
	WHERE email = $1`

	var user domain.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashPassword,
		&user.CreatedAt,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
