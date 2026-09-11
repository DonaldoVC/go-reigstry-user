package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"startplus.com/test/internal/repository"
	"startplus.com/test/internal/service"
)

func testHandler() http.Handler {
	users := service.NewUserService(repository.NewMemoryUserRepository(), []byte("test-secret"), time.Hour)
	return NewHandler(users).Routes()
}

func request(t *testing.T, handler http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func TestRegisterMissingPassword(t *testing.T) {
	response := request(t, testHandler(), "/usuarios/registro", `{"correo":"prueba@gmail.com"}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; se esperaba %d", response.Code, http.StatusBadRequest)
	}
	var payload errorResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Error != "Falta el campo contraseña" || payload.Code != "CAMPO_FALTANTE" {
		t.Fatalf("respuesta inesperada: %+v", payload)
	}
}

func TestRegisterDuplicateEmailAndPhone(t *testing.T) {
	handler := testHandler()
	valid := `{"correo":"prueba@gmail.com","telefono":"5512345678","contraseña":"Abc1@x"}`
	if got := request(t, handler, "/usuarios/registro", valid).Code; got != http.StatusCreated {
		t.Fatalf("primer registro status = %d", got)
	}
	if got := request(t, handler, "/usuarios/registro", valid).Code; got != http.StatusConflict {
		t.Fatalf("registro duplicado status = %d", got)
	}
	duplicatePhone := `{"correo":"otro@gmail.com","telefono":"55 1234 5678","contraseña":"Abc1@x"}`
	if got := request(t, handler, "/usuarios/registro", duplicatePhone).Code; got != http.StatusConflict {
		t.Fatalf("teléfono duplicado status = %d", got)
	}
}

func TestLoginReturnsJWTAndDate(t *testing.T) {
	handler := testHandler()
	register := `{"correo":"prueba@gmail.com","telefono":"5512345678","contraseña":"Abc1@x"}`
	request(t, handler, "/usuarios/registro", register)
	response := request(t, handler, "/usuarios/login", `{"correo":"prueba@gmail.com","contraseña":"Abc1@x"}`)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", response.Code, response.Body.String())
	}
	var payload map[string]string
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload["token"] == "" || payload["fecha_inicio_sesion"] == "" {
		t.Fatalf("respuesta incompleta: %+v", payload)
	}
}

func TestInvalidJSONAndUnknownField(t *testing.T) {
	if got := request(t, testHandler(), "/usuarios/login", `{`).Code; got != http.StatusBadRequest {
		t.Fatalf("JSON mal formado status = %d", got)
	}
	if got := request(t, testHandler(), "/usuarios/login", `{"correo":"a@b.com","contraseña":"Abc1@x","otro":1}`).Code; got != http.StatusBadRequest {
		t.Fatalf("campo desconocido status = %d", got)
	}
}
