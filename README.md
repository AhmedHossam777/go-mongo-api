# Go-MongoDB Course API

A RESTful API for managing courses, users, and authentication built with Go and MongoDB. This project follows a clean architecture pattern and includes JWT-based authentication, role-based access control, and Swagger documentation.

## 🚀 Features

- **User Management**: Registration, login, and user profile management.
- **Authentication**: JWT-based authentication with Access and Refresh tokens.
- **Role-Based Access Control**: Different permissions for `admin` and `user` roles.
- **Course Management**: CRUD operations for courses (Create, Read, Update, Delete).
- **Swagger Documentation**: Interactive API documentation.
- **Docker Support**: Easy deployment using Docker and Docker Compose.
- **Input Validation**: Robust request validation using `go-playground/validator`.

## 🛠️ Technologies Used

- **Go**: Version 1.25.4
- **MongoDB**: Primary database for data storage.
- **JWT**: For secure authentication and session management.
- **Swagger**: For API documentation.
- **Docker**: For containerization and deployment.
- **Go Modules**: For dependency management.

## 📋 Prerequisites

- [Go](https://go.dev/doc/install) (v1.25.4 or higher)
- [MongoDB](https://www.mongodb.com/docs/manual/installation/) (Local or Atlas)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (Optional)

## ⚙️ Configuration

Create a `.env` file in the root directory and configure the following variables:

```env
# Database Configuration
MONGODB_URI=your_mongodb_connection_string
DATABASE_NAME=coursesdb

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-key
JWT_REFRESH_SECRET=another-super-secret-key
ACCESS_TOKEN_EXPIRY_MINUTES=15
REFRESH_TOKEN_EXPIRY_DAYS=7
```

## 🏃 Running the Application

### Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/AhmedHossam777/go-mongo-api.git
   cd go-mongo-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### Using Docker

1. Build and start the containers:
   ```bash
   docker-compose up --build
   ```

The server will be available at `http://localhost:8080`.

## 📖 API Documentation

Interactive Swagger documentation is available at:
`http://localhost:8080/swagger/index.html`

## 🛣️ API Endpoints Summary

### Auth Endpoints
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and receive tokens
- `POST /api/v1/auth/refresh-tokens` - Refresh access token using refresh token
- `POST /api/v1/auth/logout` - Logout user
- `GET /api/v1/auth/active-sessions` - Get active sessions (Requires Auth)

### Course Endpoints
- `GET /api/v1/courses` - List all courses
- `GET /api/v1/courses/{id}` - Get course by ID
- `POST /api/v1/courses` - Create a new course (Requires Auth)
- `PATCH /api/v1/courses/{id}` - Update a course (Requires Auth)
- `DELETE /api/v1/courses/{id}` - Delete a course (Requires Auth)
- `DELETE /api/v1/courses/drop` - Drop all courses (Requires Admin)

### User Endpoints
- `POST /api/v1/users` - Create a user
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/me` - Get current user profile (Requires Auth)
- `GET /api/v1/users/{id}` - Get user by ID
- `PATCH /api/v1/users/{id}` - Update user details
- `DELETE /api/v1/users/{id}` - Delete user (Requires Admin)
- `DELETE /api/v1/users/drop` - Drop users collection (Requires Admin)

### General
- `GET /health` - Health check
- `GET /` - API Welcome message

## 🏗️ Project Structure

- `cmd/api/`: Entry point of the application.
- `docs/`: Swagger documentation files.
- `internal/`:
  - `config/`: Configuration and database connection logic.
  - `dto/`: Data Transfer Objects for request/response bodies.
  - `handlers/`: HTTP request handlers.
  - `helpers/`: Utility functions (JWT, password hashing, etc.).
  - `models/`: Database models.
  - `repository/`: Data access layer.
  - `services/`: Business logic layer.
- `middlewares/`: Custom HTTP middlewares (Auth, CORS, Role-based access).
- `routes/`: API route definitions.

