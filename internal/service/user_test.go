package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"startplus.com/test/internal/domain"
	"startplus.com/test/internal/repository"
)

func TestRegisterAndLogin(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	service := NewUserService(repo, []byte("test-secret"), time.Hour)
	fixedTime := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return fixedTime }

	input := RegisterInput{Email: "User@Example.com", Phone: "+52 55 1234 5678", Password: "Abc1@x"}
	if err := service.Register(context.Background(), input); err != nil {
		t.Fatalf("Register() devolvió error: %v", err)
	}

	result, err := service.Login(context.Background(), LoginInput{Email: "user@example.com", Password: "Abc1@x"})
	if err != nil {
		t.Fatalf("Login() devolvió error: %v", err)
	}
	if !result.StartedAt.Equal(fixedTime) {
		t.Fatalf("StartedAt = %v; se esperaba %v", result.StartedAt, fixedTime)
	}
	parsed, err := jwt.Parse(
		result.Token,
		func(token *jwt.Token) (any, error) { return []byte("test-secret"), nil },
		jwt.WithTimeFunc(func() time.Time { return fixedTime }),
	)
	if err != nil || !parsed.Valid {
		t.Fatalf("token JWT inválido: %v", err)
	}
}

func TestRegisterRejectsDuplicates(t *testing.T) {
	service := NewUserService(repository.NewMemoryUserRepository(), []byte("secret"), time.Hour)
	first := RegisterInput{Email: "one@example.com", Phone: "5512345678", Password: "Abc1@x"}
	if err := service.Register(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := service.Register(context.Background(), RegisterInput{Email: "ONE@example.com", Phone: "5587654321", Password: "Abc1@x"}); !errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		t.Fatalf("se esperaba correo duplicado; se obtuvo %v", err)
	}
	if err := service.Register(context.Background(), RegisterInput{Email: "two@example.com", Phone: "55 1234-5678", Password: "Abc1@x"}); !errors.Is(err, domain.ErrPhoneAlreadyRegistered) {
		t.Fatalf("se esperaba teléfono duplicado; se obtuvo %v", err)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	service := NewUserService(repository.NewMemoryUserRepository(), []byte("secret"), time.Hour)
	_ = service.Register(context.Background(), RegisterInput{Email: "one@example.com", Phone: "5512345678", Password: "Abc1@x"})
	_, err := service.Login(context.Background(), LoginInput{Email: "one@example.com", Password: "Wrong1@"})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("se esperaban credenciales inválidas; se obtuvo %v", err)
	}
}
