package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"microservices-go/services/book/internal/domain"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type bookRepo struct {
	db *pgxpool.Pool
}

type Book interface {
	Create(ctx context.Context, book *domain.Book) error
	GetByID(ctx context.Context, id int64) (*domain.Book, error)
	GetAll(ctx context.Context, title string, genres []string, filters Filters) ([]*domain.Book, error)
	Delete(ctx context.Context, id int64) error
}

func NewBookRepo(db *pgxpool.Pool) *bookRepo {
	return &bookRepo{db: db}
}

func (s *bookRepo) Create(ctx context.Context, book *domain.Book) error {
	query := `
		INSERT INTO books (title, year, genres)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []interface{}{book.Title, book.Year, book.Genres}

	return s.db.QueryRow(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (s *bookRepo) GetByID(ctx context.Context, id int64) (*domain.Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, genres, version
		FROM books
		WHERE id = $1`

	var book domain.Book

	err := s.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Year,
		&book.Genres,
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

func (s *bookRepo) GetAll(ctx context.Context, title string, genres []string, filters Filters) ([]*domain.Book, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, year, genres, version
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []any{title, genres, filters.limit(), filters.offset()}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	books := []*domain.Book{}

	for rows.Next() {
		var book domain.Book

		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Year,
			&book.Genres,
			&book.Version,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (s *bookRepo) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM books
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}
