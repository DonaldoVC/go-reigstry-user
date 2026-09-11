// Package repository supplies persistence implementations.
package repository

import (
	"context"
	"sync"

	"startplus.com/test/internal/domain"
)

// UserRepository defines the storage required by the user service.
type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	ByEmail(ctx context.Context, email string) (domain.User, bool, error)
}

// MemoryUserRepository stores users in process memory.
type MemoryUserRepository struct {
	mu      sync.RWMutex
	byEmail map[string]domain.User
	byPhone map[string]domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		byEmail: make(map[string]domain.User),
		byPhone: make(map[string]domain.User),
	}
}

// Create atomically checks uniqueness and stores a user
func (r *MemoryUserRepository) Create(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[user.Email]; exists {
		return domain.ErrEmailAlreadyRegistered
	}
	if _, exists := r.byPhone[user.Phone]; exists {
		return domain.ErrPhoneAlreadyRegistered
	}
	r.byEmail[user.Email] = user
	r.byPhone[user.Phone] = user
	return nil
}

func (r *MemoryUserRepository) ByEmail(_ context.Context, email string) (domain.User, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.byEmail[email]
	return user, exists, nil
}
