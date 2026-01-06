# 📋 Swagger Quick Reference

## Access Swagger UI
```
http://localhost:3000/swagger/index.html
```

## Regenerate Docs After Changes
```bash
~/go/bin/swag init -g cmd/api/main.go --parseDependency --parseInternal
```

## Test with Authentication

### 1. Register/Login
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Copy the accessToken from response

### 3. Use in Swagger UI
- Click "Authorize" button (top right)
- Enter: `Bearer YOUR_ACCESS_TOKEN_HERE`
- Click "Authorize" then "Close"

### 4. Use in curl
```bash
curl -X GET http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

## Endpoint Categories

### 🌐 General
- `GET /` - Home
- `GET /health` - Health check

### 🔐 Auth (`/api/v1/auth`)
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh-tokens`
- `POST /auth/logout`
- `GET /auth/active-sessions` 🔒

### 👤 Users (`/api/v1/users`)
- `GET /users?page=1&page_size=10`
- `POST /users`
- `GET /users/{id}`
- `PATCH /users/{id}`
- `DELETE /users/{id}` 🔒👑
- `GET /users/me` 🔒
- `DELETE /users/drop` 🔒👑

### 📚 Courses (`/api/v1/courses`)
- `GET /courses?page=1&page_size=10`
- `POST /courses` 🔒
- `GET /courses/{id}`
- `PATCH /courses/{id}` 🔒
- `DELETE /courses/{id}` 🔒
- `DELETE /courses/drop` 🔒👑

🔒 = Authentication Required  
👑 = Admin Role Required

## Files Structure
```
docs/
  ├── docs.go          # Generated Go package
  ├── swagger.json     # OpenAPI JSON spec
  └── swagger.yaml     # OpenAPI YAML spec
```

## Common Swagger Annotations

```go
// Main API Info (in main.go)
// @title API Title
// @version 1.0
// @description API Description
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// Endpoint (in handler)
// @Summary Short description
// @Description Detailed description  
// @Tags tag-name
// @Accept json
// @Produce json
// @Param name location type required "description"
// @Success 200 {object} Type
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /path [method]
```

## Parameter Locations
- `path` - URL path parameter (e.g., `/users/{id}`)
- `query` - URL query parameter (e.g., `?page=1`)
- `header` - HTTP header
- `body` - Request body

## Need Help?
- See `SWAGGER_SETUP_COMPLETE.md` for full guide
- See `SWAGGER.md` for detailed documentation
- Visit: https://github.com/swaggo/swag

