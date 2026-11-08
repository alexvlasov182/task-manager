# Task Manager REST API (Go & PostgreSQL)

A production-ready RESTful API service for task management, built with **Go** and **PostgreSQL**, fully containerized using **Docker Compose** and following modern security (Security Hardening) and DevOps practices.

---

## Technology Stack

- **Language:** Go (1.27)
- **Database:** PostgreSQL 16 (Alpine)
- **Containerization:** Docker & Docker Compose (Multi-stage builds)
- **Architecture:** Modular REST API, 12-Factor App methodology
- **Security:** Non-root container execution, network isolation (`127.0.0.1`), strict environment variable management.

---

## Project Architecture

```text
 ┌──────────────┐         Docker Network         ┌──────────────────┐
 │               │        (Internal Connection)   │                    │
 │  Task API     │────────────────────────────▶│    PostgreSQL      │
 │  (Container)  │    postgres:5432 / TCP         │    (Container)     │
 │  Non-root     │                                │    Healthcheck     │
 └──────┬───────┘                                └────────┬─────────┘
        │                                                   │
        │ Host Port Forwarding                              │ Volume Storage
        ▼                                                   ▼
   localhost:8080                                  postgres_data
```

---

## Quick Start (Local Setup)

### Prerequisites

Make sure the following are installed on your machine:

- Docker and Docker Compose

### 1. Clone the Repository

```bash
git clone https://github.com/oleksandr-vlasov/task-manager.git
cd task-manager
```

### 2. Configure Environment Variables

Create a `.env` file in the project root based on the example:

```bash
cp .env.example .env
```

(Or create `.env` manually with the following content):

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=supersecretpassword123
POSTGRES_DB=tasks_db
DATABASE_URL=postgres://postgres:supersecretpassword123@postgres:5432/tasks_db?sslmode=disable
```

### 3. Run with Docker Compose

Build and start the services in detached mode:

```bash
docker compose up --build -d
```

Thanks to the built-in healthcheck, the API will automatically wait for the database to be fully ready, avoiding unnecessary connection errors.

---

## API Endpoints & Usage

You can test the service using `curl`:

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Create a New Task (POST)

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Docker & DevOps", "description":"Implement multi-stage builds and security hardening", "completed":false}'
```

### 3. Get All Tasks (GET)

```bash
curl http://localhost:8080/tasks
```

---

## Implemented Security & DevOps Practices

- **Multi-stage Docker build:** The final image is based on lightweight Alpine Linux and contains only the compiled binary (~33MB).
- **Non-root user:** The application inside the container runs as an unprivileged user (`appuser`), eliminating risks in the event of a potential compromise.
- **Network isolation:** The PostgreSQL port (5432) is mapped exclusively to the local interface (`127.0.0.1`), preventing direct access from external networks.
- **Healthcheck orchestration:** Using `condition: service_healthy` guarantees a perfect startup sequence for services.

---

## License

This project was created for educational and demonstration purposes.