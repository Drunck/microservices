package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"microservices-go/services/order/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	//ErrDuplicate      = errors.New("record duplication")
	//ErrEditConflict   = errors.New("edit conflict")
)

type orderRepo struct {
	db *pgxpool.Pool
}

type Order interface {
	Insert(ctx context.Context, order *domain.Order) error
	GetByEmail(ctx context.Context, email *string) ([]*domain.Order, error)
}

func NewOrderRepo(db *pgxpool.Pool) *orderRepo {
	return &orderRepo{db: db}
}

func (s *orderRepo) Insert(ctx context.Context, order *domain.Order) error {
	query := `
	INSERT INTO orders (book_id, email)
	VALUES ($1, $2)
	RETURNING id, created_at`

	args := []any{order.BookID, order.Email}

	err := s.db.QueryRow(ctx, query, args...).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *orderRepo) GetByEmail(ctx context.Context, email *string) ([]*domain.Order, error) {
	query := `
	SELECT id, book_id, email, created_at
	FROM orders
	WHERE email = $1`
	rows, err := s.db.Query(ctx, query, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.BookID,
			&order.Email,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()

	return orders, nil
}
