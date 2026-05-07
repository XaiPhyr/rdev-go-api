# Go Backend Engine

A high-performance, containerized REST API built with **Golang 1.26**. This project implements a robust backend architecture focused on security, system performance, and automated deployment workflows.

## 🚀 Technical Highlights

*   **SQL-First Design:** Utilizes **Bun ORM** with **PostgreSQL** for type-safe, high-performance database interactions without the overhead of heavy traditional ORMs.
*   **Security Architecture:** Implements a custom **Role-Based Access Control (RBAC)** system with hierarchical permissions and junction tables.
*   **Performance Optimization:** Integrated **Redis** for authorization caching to reduce database latency.
*   **Infrastructure as Code:** Fully containerized using **Docker** and **Docker Compose** with multi-stage builds to ensure a minimal production footprint.
*   **CI/CD Automation:** A unified **Makefile** manages the lifecycle from local development to automated **SSH deployment** via **GitHub Actions**.

---

## 🛠 Getting Started

### Prerequisites
*   Go 1.26+
*   Docker & Docker Compose
*   Make (optional, but recommended)

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

---