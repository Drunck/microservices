package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"microservices-go/pkg/hash"
	"microservices-go/pkg/token"
	"microservices-go/pkg/validator"
	"microservices-go/services/user/internal/domain"
	"microservices-go/services/user/internal/repository"
	"time"
)

var (
	ErrFailedValidation = errors.New("validation failed")
	ErrWrongCredentials = errors.New("wrong user credentials")
	ErrDuplicate        = errors.New("record duplication")
)

type UserSignUpDTO struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type UserSignInDTO struct {
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type Token struct {
	PlainText string
}

type UserService interface {
	SignUp(ctx context.Context, user UserSignUpDTO) error
	SignIn(ctx context.Context, user UserSignInDTO) (Token, error)
}

type service struct {
	repo         repository.User
	hasher       hash.PasswordHasher
	tokenManager token.TokenManager
}

func New(repo repository.User, hasher hash.PasswordHasher, tokenManager token.TokenManager) *service {
	return &service{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *service) SignUp(ctx context.Context, input UserSignUpDTO) error {
	passwordHash, err := s.hasher.Hash(input.HashPassword)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		HashPassword: passwordHash,
	}

	v := validator.New()
	validateUser(v, &user)
	if !v.Valid() {
		return ErrFailedValidation
	}

	err = s.repo.Insert(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicate):
			return ErrDuplicate
		default:
			return err
		}
	}
	return nil

}

func (s *service) SignIn(ctx context.Context, input UserSignInDTO) (Token, error) {
	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return Token{}, ErrWrongCredentials
		default:
			return Token{}, err
		}
	}

	err2 := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(input.HashPassword))
	if err2 != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return Token{}, nil
		default:
			return Token{}, err
		}
	}

	token, err := s.tokenManager.NewToken(user.ID, time.Duration(12)*time.Hour)
	return Token{PlainText: token}, err

}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 64, "password", "must not be more than 64 bytes long")
}

func validateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "name must be provided")
	v.Check(len(name) <= 100, "name", "name must not be more than 100 bytes long")
}

func validateUser(v *validator.Validator, user *domain.User) {
	validateEmail(v, user.Email)
	validateName(v, user.Name)
	validatePassword(v, user.HashPassword)
}
