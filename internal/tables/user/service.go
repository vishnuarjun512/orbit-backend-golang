package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"orbit-backend-golang/internal/security"
)

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidInput         = errors.New("invalid registration input")
	ErrUserNotRegistered    = errors.New("User not registered")
	ErrInvalidPassword      = errors.New("Invalid Password")
	ErrTokenGenerationError = errors.New("Token not generated")
)

type Service struct {
	repo       *Repository
	jwtService *security.JWTService
}

func NewService(repo *Repository, jwtService *security.JWTService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	// Normalize email.
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Validate input.
	if email == "" || strings.TrimSpace(req.Password) == "" {
		return ErrInvalidInput
	}

	// Hash password using Argon2id.
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Prepare user.
	newUser := &User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	// Save user.
	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(req.Email))

	//Validate Email

	if email == "" || strings.TrimSpace(req.Password) == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotRegistered
	}

	isPassValid, err := security.VerifyPassword(req.Password, user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !isPassValid {
		return nil, ErrInvalidPassword
	}

	accessToken, err := s.jwtService.CreateToken(user.UserID)

	if err != nil {
		fmt.Print(err)
		return nil, ErrTokenGenerationError
	}

	safeUser := SafeUser{
		Email:     user.Email,
		UserID:    user.UserID,
		Username:  user.Username,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	return &LoginResponse{safeUser, accessToken}, nil
}
