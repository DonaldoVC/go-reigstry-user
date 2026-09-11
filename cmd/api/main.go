// Package main initializes and starts the user HTTP API.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"startplus.com/test/internal/httpapi"
	"startplus.com/test/internal/repository"
	"startplus.com/test/internal/service"
)

func main() {
	// Load local configuration without overriding existing environment variables.
	// A missing .env is allowed when configuration is provided by the environment.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("No se pudo cargar el archivo .env; revisa su formato y permisos")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET es obligatorio")
	}

	ttl := 24 * time.Hour
	if value := os.Getenv("JWT_TTL"); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil || parsed <= 0 {
			log.Fatal("JWT_TTL debe ser una duración positiva (por ejemplo, 24h)")
		}
		ttl = parsed
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := repository.NewMemoryUserRepository()
	users := service.NewUserService(repo, []byte(secret), ttl)
	handler := httpapi.NewHandler(users)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("servicio escuchando en http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
