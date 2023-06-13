package service

import (
	"context"
	"errors"
	"microservices-go/pkg/validator"
	"microservices-go/services/book/internal/domain"
	"microservices-go/services/book/internal/repository"
	"time"
)

var (
	ErrFailedValidation = errors.New("validation failed")
	ErrDuplicate        = errors.New("record duplication")
)

type CreateBookDTO struct {
	Title  string   `json:"title"`
	Year   int32    `json:"year"`
	Genres []string `json:"genres"`
}

type BookService interface {
	CreateBook(ctx context.Context, book CreateBookDTO) error
	GetBookByID(ctx context.Context, id int64) (*domain.Book, error)
	GetBooks(ctx context.Context, title string, genres []string, filters repository.Filters) ([]*domain.Book, error)
}

type service struct {
	repo repository.Book
}

func New(repo repository.Book) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateBook(ctx context.Context, input CreateBookDTO) error {
	book := domain.Book{
		Title:  input.Title,
		Year:   input.Year,
		Genres: input.Genres,
	}

	v := validator.New()

	if ValidateBook(v, &book); !v.Valid() {
		return ErrFailedValidation
	}

	err := s.repo.Create(ctx, &book)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetBookByID(ctx context.Context, id int64) (*domain.Book, error) {
	book, err := s.repo.GetByID(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return book, nil
}

func (s *service) GetBooks(ctx context.Context, title string, genres []string, filters repository.Filters) ([]*domain.Book, error) {
	v := validator.New()

	if repository.ValidateFilters(v, filters); !v.Valid() {
		return nil, ErrFailedValidation
	}
	var books []*domain.Book

	books, err := s.repo.GetAll(ctx, title, genres, filters)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return books, err
}

func ValidateBook(v *validator.Validator, book *domain.Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year >= 1888, "year", "must be greater than 1888")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}
