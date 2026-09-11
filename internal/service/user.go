// Package service implements user registration and authentication use cases.
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"startplus.com/test/internal/domain"
	"startplus.com/test/internal/repository"
	"startplus.com/test/internal/validation"
)

var (
	ErrInvalidEmail    = errors.New("el correo no es válido")
	ErrInvalidPhone    = errors.New("el teléfono no es válido")
	ErrInvalidPassword = errors.New("la contraseña debe tener entre 6 y 12 caracteres e incluir una mayúscula, una minúscula, un número y uno de estos caracteres: @, $ o &")
)

type RegisterInput struct {
	Email    string
	Phone    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token     string
	StartedAt time.Time
}

type UserService struct {
	repository repository.UserRepository
	jwtSecret  []byte
	tokenTTL   time.Duration
	now        func() time.Time
}

func NewUserService(repo repository.UserRepository, jwtSecret []byte, tokenTTL time.Duration) *UserService {
	return &UserService{repository: repo, jwtSecret: jwtSecret, tokenTTL: tokenTTL, now: time.Now}
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) error {
	email := validation.NormalizeEmail(input.Email)
	phone := validation.NormalizePhone(input.Phone)
	if !validation.ValidEmail(email) {
		return ErrInvalidEmail
	}
	if !validation.ValidPhone(phone) {
		return ErrInvalidPhone
	}
	if !validation.ValidPassword(input.Password) {
		return ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("crear hash de contraseña: %w", err)
	}
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return fmt.Errorf("crear identificador de usuario: %w", err)
	}

	return s.repository.Create(ctx, domain.User{
		ID:           hex.EncodeToString(idBytes),
		Email:        email,
		Phone:        phone,
		PasswordHash: hash,
		CreatedAt:    s.now().UTC(),
	})
}

func (s *UserService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	email := validation.NormalizeEmail(input.Email)
	if !validation.ValidEmail(email) {
		return LoginResult{}, ErrInvalidEmail
	}
	user, exists, err := s.repository.ByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, fmt.Errorf("consultar usuario: %w", err)
	}
	if !exists || bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password)) != nil {
		return LoginResult{}, domain.ErrInvalidCredentials
	}

	startedAt := s.now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(startedAt),
		ExpiresAt: jwt.NewNumericDate(startedAt.Add(s.tokenTTL)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return LoginResult{}, fmt.Errorf("firmar token: %w", err)
	}
	return LoginResult{Token: token, StartedAt: startedAt}, nil
}
