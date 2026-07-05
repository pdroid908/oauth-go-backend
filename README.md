# 🔐 OAuth Go Backend

A secure and scalable OAuth authentication backend built with **Go**. This project provides authentication services, OAuth integration, JWT-based authorization, session management, and RESTful APIs for modern web applications.

Designed with clean architecture principles, this backend is intended to work alongside a separate frontend application while keeping authentication and sensitive credentials securely managed on the server.

> **Status:** 🚧 Under Development

---

# ✨ Features

- 🔑 OAuth Authentication
- 🔐 JWT Access & Refresh Tokens
- 👤 User Authentication
- 🛡️ Secure Password Hashing
- 🔄 Token Refresh Flow
- 🚪 Logout & Session Invalidation
- 📦 RESTful API
- ⚡ High Performance with Go
- 🗄 PostgreSQL Integration
- ⚙️ Environment-based Configuration

---

# 🛠 Tech Stack

| Technology | Purpose |
|------------|---------|
| Go | Backend Language |
| Gin / Fiber |
| PostgreSQL | Database |
| Cache & Session Storage |
| JWT | Authentication |
| OAuth 2.0 | Authorization |
| Docker | Containerization |
| Docker Compose | Local Development |

---

# 🚀 Getting Started

## Clone the repository

```bash
git clone https://github.com/pdroid908/oauth-go-backend.git
```

Move into the project

```bash
cd oauth-go-backend
```

Install dependencies

```bash
go mod tidy
```

Run the application

```bash
go run cmd/main.go
```

---

# 🐳 Docker

Start all services

```bash
docker compose up --build
```

Stop services

```bash
docker compose down
```

---

# 🔐 Authentication Flow

```
User
 │
 ▼
Frontend
 │
 ▼
OAuth Provider
 │
 ▼
OAuth Go Backend
 │
 ├── Verify OAuth Identity
 ├── Create or Find User
 ├── Generate JWT
 ├── Store Refresh Token
 └── Return Authentication Response
 │
 ▼
Frontend
```

---

# 📡 API Endpoints

| Method | Endpoint | Description |
|---------|----------|-------------|
| GET | /auth/google | Google OAuth Login |
| GET | /auth/github | GitHub OAuth Login |
| GET | /auth/callback | OAuth Callback |
| POST | /auth/login | Email Login |
| POST | /auth/register | Register |
| POST | /auth/refresh | Refresh Access Token |
| POST | /auth/logout | Logout |
| GET | /user/profile | Current User |

---

# 🔒 Security

This project follows common backend security practices, including:

- Password hashing using bcrypt
- JWT Authentication
- Refresh Token rotation
- Environment variable configuration
- Protected API routes
- Input validation
- SQL Injection prevention through parameterized queries
- Secure OAuth flow
- CORS configuration
- HTTP security headers

---

# 📈 Roadmap

- [ ] Google OAuth
- [ ] GitHub OAuth
- [ ] Email & Password Login
- [ ] JWT Authentication
- [ ] Refresh Tokens
- [ ] Role-Based Access Control (RBAC)
- [ ] Email Verification
- [ ] Password Reset
- [ ] Docker Deployment
- [ ] Swagger Documentation
- [ ] Unit Testing
- [ ] Integration Testing
- [ ] Rate Limiting
- [ ] Audit Logging

---

# 🧪 Testing

Run all tests

```bash
go test ./...
```

---

# 🌐 Related Projects

| Project | Description |
|---------|-------------|
| Frontend OAuth | Next.js frontend client |
| OAuth Go Backend | Authentication API |

---

# 👨‍💻 Author

**Putra Rohman**

Backend Developer passionate about building secure, scalable, and maintainable backend systems.

### Core Skills

- Go
- REST API
- OAuth 2.0
- JWT
- PostgreSQL
- Docker
- TypeScript
- Next.js
- Clean Architecture

GitHub

https://github.com/pdroid908

---

# ⭐ Support

If you find this project useful, consider giving it a ⭐ on GitHub.

---

## License

Licensed under the MIT License.
