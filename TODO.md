# TODO: Project Enhancements & Future Scope

This document outlines the planned improvements and features to elevate this project for production-readiness and professional portfolio standards.

### 🛡️ Security
- [ ] **Refresh Token Implementation**: Add support for refresh tokens to improve security and user experience.
- [ ] **Rate Limiting**: Implement middleware to prevent brute-force attacks and API abuse.
- [ ] **Secure Headers**: Integrate `helmet`-like middleware for Go (e.g., using `chi/middleware` or manual headers) to set security-related HTTP headers.
- [ ] **Input Sanitization**: Add deeper sanitization for user-generated content to prevent XSS or injection attacks.
- [ ] **Environment Variable Validation**: Enhance `config/database.go` to strictly validate required environment variables at startup.

### 🚀 Performance & Scalability
- [ ] **Caching**: Integrate Redis for caching frequently accessed data like course lists or user profiles.
- [ ] **Database Indexing**: Define and implement MongoDB indexes for `email` in users and `title/category` in courses to optimize query performance.
- [ ] **Pagination**: Add pagination, filtering, and sorting for `GetAllUsers` and `GetAllCourses` endpoints.
- [ ] **Connection Pooling**: Fine-tune MongoDB client options for better connection management under load.

### 🧪 Testing & Quality
- [ ] **Unit Testing**: Increase coverage for `services` and `helpers` using `testify/assert`.
- [ ] **Integration Testing**: Add integration tests for API endpoints using `net/http/httptest`.
- [ ] **Mocking**: Implement a mocking layer for repositories to test services in isolation.
- [ ] **CI/CD Pipeline**: Set up GitHub Actions for automated linting (`golangci-lint`) and testing on every push.

### 📖 Documentation & DX (Developer Experience)
- [ ] **Swagger/OpenAPI**: Integrate `swaggo/swag` to automatically generate interactive API documentation.
- [ ] **Graceful Shutdown**: Implement graceful shutdown in `main.go` to handle `SIGTERM` and `SIGINT` signals properly.
- [ ] **Structured Logging**: Replace standard `log` with a structured logger like `zap` or `zerolog` for better observability.
- [ ] **Custom Error Types**: Refactor error handling to use custom internal error types for better status code mapping.

### ✨ Features
- [ ] **File Uploads**: Add support for uploading course images or user avatars (using AWS S3 or Local Storage).
- [ ] **Email Service**: Implement an email service for password resets and welcome emails.
- [ ] **Search**: Add a search endpoint with fuzzy matching for courses.
