# 🔐 OAuth Go Backend

A secure authentication service built with **Go** that provides OAuth authentication, session management, and REST APIs for modern web applications.

The backend authenticates users, stores sessions securely using **HTTP-only cookies**, and exposes protected endpoints for authenticated clients.

## Features

- OAuth Authentication
- Email & Password Authentication
- JWT Authentication
- HTTP-only Cookie Sessions
- Refresh Token
- User Registration
- Protected APIs
- PostgreSQL Integration
- Redis Session Storage
- RESTful API

## Tech Stack

- Go
- PostgreSQL
- JWT
- OAuth 2.0
- Docker

## Authentication Flow

```text
Frontend
    │
    ▼
OAuth Go Backend
    │
    ├── Validate Credentials
    ├── Generate JWT
    ├── Store Refresh Token
    ├── Set HTTP-only Cookie
    ▼
Authenticated User
```

## Security

- HTTP-only Cookies
- Secure Cookie Configuration
- Password Hashing (bcrypt)
- JWT Authentication
- Refresh Token Rotation
- Input Validation
- SQL Injection Protection
- CORS Configuration

## API

| Method | Endpoint |
|---------|----------|
| POST | /login |
| POST | /register |
| POST | /logout |
| POST | /refresh |
| GET | /me |

## Getting Started

```bash
git clone https://github.com/pdroid908/oauth-go-backend.git

cd oauth-go-backend

go mod tidy

go run main.go
```

## Environment

```env
PORT=8080

DATABASE_URL=

JWT_SECRET=

REDIS_URL=
```

## Frontend

Frontend repository:

https://github.com/pdroid908/frontend-oauth

## Author

**Putra Rohman**

Backend Developer specializing in Go, TypeScript, PostgreSQL, Redis, Docker, and REST APIs.
