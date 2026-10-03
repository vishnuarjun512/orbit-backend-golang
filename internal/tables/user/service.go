package user

import (
	"context"
	"errors"
	"strings"

	"orbit-backend-golang/internal/security"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidInput          = errors.New("invalid registration input")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// Normalize user input.
	email := strings.TrimSpace(strings.ToLower(req.Email))
	username := strings.TrimSpace(strings.ToLower(req.Username))
	fullName := strings.TrimSpace(req.FullName)

	if email == "" || username == "" || fullName == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	// Check whether the email is already registered.
	existingEmail, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingEmail != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check whether the username is already taken.
	existingUsername, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if existingUsername != nil {
		return nil, ErrUsernameAlreadyExists
	}

	// Hash the password using Argon2id.
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Prepare the user model.
	newUser := &User{
		Email:        email,
		Username:     username,
		FullName:     fullName,
		PasswordHash: passwordHash,
	}

	// Save the user.
	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
