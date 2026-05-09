# Go Backend Engine

A high-performance, containerized REST API built with **Golang 1.26**. This project implements a robust backend architecture focused on security, system performance, and automated deployment workflows.

## 🚀 Technical Highlights

- **SQL-First Design:** Utilizes **Bun ORM** with **PostgreSQL** for type-safe, high-performance database interactions without the overhead of heavy traditional ORMs.
- **Security Architecture:** Implements a custom **Role-Based Access Control (RBAC)** system with hierarchical permissions and junction tables.
- **Performance Optimization:** Integrated **Redis** for authorization caching to reduce database latency.
- **Infrastructure as Code:** Fully containerized using **Docker** and **Docker Compose** with multi-stage builds to ensure a minimal production footprint.
- **CI/CD Automation:** A unified **Makefile** manages the lifecycle from local development to automated **SSH deployment** via **GitHub Actions**.

## 🛠 Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Make (optional, but recommended)

### Local Development (Docker)

The environment is configured to bridge the containerized API with your local database via `host-gateway`.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/XaiPhyr/rdev-go-api.git
   cd rdev-go-api
   ```

2. **Spin up the environment:**
   ```bash
   docker compose up -d --build
   ```
   The API will be accessible at `http://localhost:8200/api/v1`.

## 📂 Folder structure
```
.
├── cmd/
│   ├── api/                # Main entry point; initializes the server, services, and database
│   └── migrate/            # CLI tool for managing database schema versions
├── internal/
│   ├── config/             # Configuration management; loads environment variables
│   ├── data/               # The "Source of Truth"; contains database models and repositories
│   │   └── migrations/     # SQL files defining schema changes over time
│   ├── dto/                # Data Transfer Objects; handles request/response shapes and query sanitization
│   ├── server/             # HTTP layer; contains route definitions and controller handlers
│   ├── service/            # Business logic layer; bridges DTOs and Data models
│   └── templates/          # HTML templates for email services
├── go.mod                  # Project dependencies
└── go.sum                  # Dependency checksums
```

## Import Flow: _Server_ -> _Service_ -> _Data_

```
| Layer      | Imports DTO? | Imports Data? | Imports Service? |
| ---------- | ------------ | ------------- | ---------------- |
| Handler    | Yes          | No            | Yes              |
| Service    | Yes          | Yes           | No               |
| Repository | No*          | Yes           | No               |
*except for dto.BaseFilters. Service and Repository can see it without causing a cycle error
```
