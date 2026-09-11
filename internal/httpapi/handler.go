// Package httpapi exposes the application through JSON over HTTP.
package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"startplus.com/test/internal/domain"
	"startplus.com/test/internal/service"
)

type Handler struct {
	users *service.UserService
}

func NewHandler(users *service.UserService) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /usuarios/registro", h.register)
	mux.HandleFunc("POST /usuarios/login", h.login)
	mux.HandleFunc("GET /salud", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
	})
	return mux
}

type registerRequest struct {
	Email    *string `json:"correo"`
	Phone    *string `json:"telefono"`
	Password *string `json:"contraseña"`
}

type loginRequest struct {
	Email    *string `json:"correo"`
	Password *string `json:"contraseña"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"codigo"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if request.Email == nil || strings.TrimSpace(*request.Email) == "" {
		writeError(w, http.StatusBadRequest, "Falta el campo correo", "CAMPO_FALTANTE")
		return
	}
	if request.Password == nil || *request.Password == "" {
		writeError(w, http.StatusBadRequest, "Falta el campo contraseña", "CAMPO_FALTANTE")
		return
	}
	if request.Phone == nil || strings.TrimSpace(*request.Phone) == "" {
		writeError(w, http.StatusBadRequest, "Falta el campo teléfono", "CAMPO_FALTANTE")
		return
	}

	err := h.users.Register(r.Context(), service.RegisterInput{
		Email: *request.Email, Phone: *request.Phone, Password: *request.Password,
	})
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyRegistered):
		writeError(w, http.StatusConflict, "El correo ya se encuentra registrado", "CORREO_REGISTRADO")
	case errors.Is(err, domain.ErrPhoneAlreadyRegistered):
		writeError(w, http.StatusConflict, "El teléfono ya se encuentra registrado", "TELEFONO_REGISTRADO")
	case errors.Is(err, service.ErrInvalidEmail):
		writeError(w, http.StatusUnprocessableEntity, err.Error(), "CORREO_INVALIDO")
	case errors.Is(err, service.ErrInvalidPhone):
		writeError(w, http.StatusUnprocessableEntity, err.Error(), "TELEFONO_INVALIDO")
	case errors.Is(err, service.ErrInvalidPassword):
		writeError(w, http.StatusUnprocessableEntity, err.Error(), "CONTRASENA_INVALIDA")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "Error interno del servidor", "ERROR_INTERNO")
	default:
		writeJSON(w, http.StatusCreated, map[string]string{"mensaje": "Usuario registrado correctamente"})
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if request.Email == nil || strings.TrimSpace(*request.Email) == "" {
		writeError(w, http.StatusBadRequest, "Falta el campo correo", "CAMPO_FALTANTE")
		return
	}
	if request.Password == nil || *request.Password == "" {
		writeError(w, http.StatusBadRequest, "Falta el campo contraseña", "CAMPO_FALTANTE")
		return
	}

	result, err := h.users.Login(r.Context(), service.LoginInput{Email: *request.Email, Password: *request.Password})
	switch {
	case errors.Is(err, service.ErrInvalidEmail):
		writeError(w, http.StatusUnprocessableEntity, err.Error(), "CORREO_INVALIDO")
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "Correo o contraseña incorrectos", "CREDENCIALES_INVALIDAS")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "Error interno del servidor", "ERROR_INTERNO")
	default:
		writeJSON(w, http.StatusOK, map[string]string{
			"token": result.Token, "fecha_inicio_sesion": result.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	if contentType := r.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type debe ser application/json", "TIPO_CONTENIDO_INVALIDO")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "El cuerpo debe contener un objeto JSON válido", "JSON_INVALIDO")
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "El cuerpo debe contener un único objeto JSON", "JSON_INVALIDO")
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, errorResponse{Error: message, Code: code})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
