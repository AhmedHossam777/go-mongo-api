# Swagger API Documentation

This project uses Swagger (OpenAPI) for API documentation. The documentation is automatically generated from code annotations.

## Accessing Swagger UI

Once the server is running, you can access the interactive Swagger UI at:

```
http://localhost:3000/swagger/index.html
```

## API Documentation

The Swagger UI provides:
- **Interactive API testing**: Try out API endpoints directly from the browser
- **Request/Response schemas**: View the structure of requests and responses
- **Authentication**: Test authenticated endpoints using Bearer tokens
- **Complete API reference**: All endpoints, parameters, and responses documented

## API Endpoints Overview

### Authentication (`/api/v1/auth`)
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and get access tokens
- `POST /auth/refresh-tokens` - Refresh access token
- `POST /auth/logout` - Logout and invalidate token
- `GET /auth/active-sessions` - Get active sessions (requires auth)

### Users (`/api/v1/users`)
- `GET /users` - Get all users (paginated)
- `POST /users` - Create a new user
- `GET /users/{id}` - Get user by ID
- `PATCH /users/{id}` - Update user by ID
- `DELETE /users/{id}` - Delete user (admin only)
- `GET /users/me` - Get current user profile (requires auth)
- `DELETE /users/drop` - Drop all users (admin only)

### Courses (`/api/v1/courses`)
- `GET /courses` - Get all courses (paginated)
- `POST /courses` - Create a new course (requires auth)
- `GET /courses/{id}` - Get course by ID
- `PATCH /courses/{id}` - Update course (requires auth)
- `DELETE /courses/{id}` - Delete course (requires auth)
- `DELETE /courses/drop` - Drop all courses (admin only)

## Authentication in Swagger UI

To test authenticated endpoints:

1. First, login or register using the auth endpoints
2. Copy the `accessToken` from the response
3. Click the **"Authorize"** button at the top of Swagger UI
4. Enter `Bearer <your-access-token>` in the value field
5. Click **"Authorize"** and then **"Close"**
6. Now you can test protected endpoints

## Regenerating Documentation

If you modify the API or add new endpoints, regenerate the Swagger docs:

```bash
~/go/bin/swag init -g cmd/api/main.go --parseDependency --parseInternal
```

Or if swag is in your PATH:

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

## Swagger Annotations

The documentation is generated from special comments in the code:

### Main API Info (in `cmd/api/main.go`):
```go
// @title Go-MongoDB Course API
// @version 1.0
// @description A RESTful API for managing courses, users, and authentication
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
```

### Endpoint Documentation (in handlers):
```go
// @Summary Short description
// @Description Detailed description
// @Tags tag-name
// @Accept json
// @Produce json
// @Param name location type required "description"
// @Success 200 {object} ModelType
// @Router /path [method]
```

## Files Generated

The following files are auto-generated in the `docs/` directory:
- `docs.go` - Go package with embedded documentation
- `swagger.json` - OpenAPI specification in JSON format
- `swagger.yaml` - OpenAPI specification in YAML format

**Note**: Do not edit these files manually. They will be overwritten when regenerating docs.

## Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

