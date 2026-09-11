// Package domain contains the application's core entities and errors.
package domain

import (
	"errors"
	"time"
)

// User is a registered user. PasswordHash is never exposed through the API.
type User struct {
	ID           string
	Email        string
	Phone        string
	PasswordHash []byte
	CreatedAt    time.Time
}

var (
	ErrEmailAlreadyRegistered = errors.New("el correo ya se encuentra registrado")
	ErrPhoneAlreadyRegistered = errors.New("el teléfono ya se encuentra registrado")
	ErrInvalidCredentials     = errors.New("correo o contraseña incorrectos")
)
