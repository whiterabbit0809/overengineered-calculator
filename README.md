# Overengineered Calculator (Go backend)

Small over-engineered backend in Go with email/password auth, JWT-based protected endpoints, a stateful calculator, and per-user calculation history stored in PostgreSQL. Deployed as a Dockerized app on Render, with a tiny HTML UI and a Postman collection for tests.

## Features

- Email/password signup & login (bcrypt-hashed passwords)
- JWT authentication (`Authorization: Bearer <token>`)
- Stateful calculator (ADD, SUBTRACT, MULTIPLY, DIVIDE)
- Per-user calculation history in Postgres
- Minimal HTML frontend for manual testing
- Postman collection for end-to-end tests
- Unit tests for auth and calculator logic

## Tech stack

- Go 1.x (net/http, database/sql)
- PostgreSQL
- JWT (github.com/golang-jwt/jwt/v5)
- bcrypt (golang.org/x/crypto/bcrypt)
- Docker + docker-compose
- Render Web Service
- Postman, testify

## API

All endpoints are served under the `/api/v1` prefix.

- Auth endpoints:
  - `POST /api/v1/auth/signup`
  - `POST /api/v1/auth/login`
- Calculator:
  - `POST /api/v1/calc` (protected)
- History:
  - `GET /api/v1/history` (protected)

Protected endpoints require:

```http
Authorization: Bearer <jwt-token>
