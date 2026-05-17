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
```text
.
├── cmd/
│   ├── app/                # Main entry point; initializes game loop, services, and db
│   │   └── main.go         # Calls ebiten.RunGame()
│   └── migrate/            # CLI tool for managing database schema versions
├── internal/
│   ├── config/             # Configuration management; loads environment variables
│   ├── db/                 # Database connection pools & global migrations
│   │   └── migrations/     # SQL files defining schema changes over time
│   ├── shared/             # Global utilities shared across multiple domains
│   │   └── dto/
│   │       ├── query.go    # Global query params (Pagination, Sorting, Search)
│   │       └── response.go # Standardized global API response envelopes
│   ├── users/              # Self-contained User Domain
│   │   ├── handler.go      # HTTP controllers / Input routers for users
│   │   ├── repository.go   # Database SQL queries and storage operations
│   │   ├── service.go      # Business logic and user-specific validation
│   │   └── types.go        # User-specific request/response shapes (local DTOs)
│   ├── orders/             # Self-contained Order Domain
│   │   ├── handler.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── types.go
│   └── templates/          # HTML templates (e.g., for system emails)
├── go.mod                  # Go module tracking project dependencies
└── go.sum                  # System dependency checksums
```

## Package Dependency & Data Flow

In this domain-driven architecture, components within the same folder (e.g., `handler.go`, `service.go`, `repository.go`) belong to the same Go package. They reference each other's types directly without Go `import` statements. 

### Data Flow Matrix

| Component | Can Access Local Types? (`types.go`) | Can Import `shared/dto`? | Can Call Other Domains? |
| :--- | :--- | :--- | :--- |
| **Handler** | Yes (Requests/Responses) | Yes (Global Queries) | No (Routes through its own Service) |
| **Service** | Yes (Core business entities) | Yes (Unpacks Global Queries) | Only via **Interfaces** (to avoid cycles) |
| **Repository**| Yes (Database models/mappings) | No | No (Strictly handles its own domain DB) |

### Key Rules for Domain Isolation:
1. **No Direct Cross-Imports:** `internal/orders` must never directly import `internal/users`. If the orders domain needs user data, it must define a local interface that the user service satisfies, or data must be aggregated at the handler/orchestration level.
2. **Database Purity:** Repositories only talk to the database driver/pool (`internal/db`) and map data to local types. They remain completely blind to HTTP requests, global DTOs, or other feature domains.
