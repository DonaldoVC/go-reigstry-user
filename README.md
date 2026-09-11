# Servicio de usuarios en Go

API REST mínima para registrar usuarios y autenticar mediante JWT. Usa una
arquitectura por capas (`httpapi`, `service`, `repository`, `domain`) y almacena
los usuarios en memoria. Las contraseñas se guardan con bcrypt.

## Ejecutar

Requiere Go 1.23 o Docker.

Crea un archivo `.env` en la raíz del proyecto usando `.env.example` como
referencia y configura `JWT_SECRET`. Se carga automáticamente al iniciar, sin
sobrescribir variables de entorno existentes. Ejecuta desde la raíz del proyecto;
en IntelliJ configura esa carpeta como directorio de trabajo. Si no hay `.env`,
puedes proporcionar las variables directamente desde el entorno.

```bash
go mod download
go run ./cmd/api
```

Con Docker:

```bash
docker build -t user-service .
docker run --rm -p 8080:8080 -e JWT_SECRET='un-secreto-largo-y-seguro' user-service
```

## Endpoints

### Registro

`POST /usuarios/registro`

```json
{
  "correo": "prueba@gmail.com",
  "telefono": "+525512345678",
  "contraseña": "Abc1@x"
}
```

Devuelve `201 Created`. Un correo o teléfono repetido devuelve `409 Conflict`.
La contraseña debe tener de 6 a 12 caracteres y al menos una mayúscula, una
minúscula, un número y uno de `@`, `$` o `&`.

### Login

`POST /usuarios/login`

```json
{
  "correo": "prueba@gmail.com",
  "contraseña": "Abc1@x"
}
```

Respuesta `200 OK`:

```json
{
  "token": "eyJ...",
  "fecha_inicio_sesion": "2026-09-10T12:00:00Z"
}
```

Todos los errores tienen la forma:

```json
{
  "error": "Falta el campo contraseña",
  "codigo": "CAMPO_FALTANTE"
}
```

El teléfono acepta de 10 a 15 dígitos y un `+` inicial opcional. También se
normalizan espacios, guiones y paréntesis.

## Pruebas

```bash
go test ./...
```

Para ver cobertura:

```bash
go test -cover ./...
```

## Consideración de producción

El repositorio en memoria pierde sus datos al reiniciar. La interfaz
`UserRepository` permite sustituirlo por PostgreSQL u otra base de datos sin
cambiar las reglas de negocio ni los handlers HTTP.
