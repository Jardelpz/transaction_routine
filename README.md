# Transaction Routine

Simple REST API to manage accounts and transactions.


---

# Features

* Create account
* Retrieve account by ID
* Create transaction
* PostgreSQL database
* Audit logging
* Structured logging
* Swagger / OpenAPI documentation
* Docker support
* Makefile automation

---

# Tech Stack

* Go
* Gin
* PostgreSQL
* Docker
* Swagger / OpenAPI
* Makefile

---

# Running the Project

Make sure Docker is installed and running.

## 1. Create environment file

Copy the sample environment file:

```bash
cp .env-sample .env
```

---

## 2. Start the application

Using Makefile:

```bash
make up
```

Or manually:

```bash
docker compose up --build
```

---

# API Base URL

```text
http://localhost:8080
```

---

# Swagger Documentation

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

OpenAPI file:

```text
http://localhost:8080/openapi.yaml
```

---

# Makefile Commands

Common commands:

```bash
make up        # start application with Docker
make down      # stop containers
make reset     # build application
make test      # run tests
```

---

# Project Structure

```text
cmd/                # application entrypoint
internal/           # core business logic
scripts/            # database initialization scripts
docker-compose.yml  # container orchestration
Dockerfile          # application container
```
