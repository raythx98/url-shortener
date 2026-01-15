# Gemini Project Context

## Overview
This project is a URL shortener service written in Go.
It currently uses PostgreSQL for storage and runs on AWS (EC2) with Nginx as a reverse proxy.
The goal is to migrate this project to Oracle Cloud Infrastructure (OCI).

## Tech Stack
- **Language:** Go 1.24.0
- **Database:** PostgreSQL
- **Infrastructure:** Currently AWS (EC2), migrating to OCI.
- **Tools:** Docker, Make, SQLC, Swagger, Migrate.

## Migration Analysis (Target: `../go-dutch`)

| Feature | Current (`url-shortener`) | Target (`go-dutch`) | Planned Action |
| :--- | :--- | :--- | :--- |
| **Infrastructure** | AWS EC2 (Amazon Linux/CentOS-like) | OCI (Oracle Linux, Podman) | Update workflow to support OCI/Podman. |
| **Reverse Proxy** | Nginx (Manual Install on Host) | Caddy (Dockerized Service) | Remove Nginx. Add Caddy to Docker Compose. |
| **Orchestration** | Manual `docker run` commands in Makefile | Docker Compose | Create `docker-compose.yaml` to manage App, DB, Migration, and Caddy. |
| **Container Image** | `debian:bookworm` (Root user default) | `distroless/static-debian12:nonroot` | Update `Dockerfile` to use distroless and non-root user. |
| **App Port** | 5051 | 8080 | Update `Dockerfile` and Env vars to use port 8080. |
| **Go Version** | 1.24.0 | 1.24.11 | Update `go.mod`, `Dockerfile`, and CI/CD. |
| **CI/CD** | `.github/workflows/aws.yml` (SSH + Shell) | `.github/workflows/deploy.yml` (SSH + Podman setup) | Rename `aws.yml` to `deploy.yml` and adopt `go-dutch` deployment steps. |

## Migration Plan

### 1. Cleanup
- [ ] Remove `.nginx/` directory.

### 2. Configuration Files
- [ ] **Create `docker-compose.yaml`**:
    - Define services: `caddy`, `app`, `migrate`, `db`.
    - Align volume and network configuration with `go-dutch`.
    - Use `APP_URLSHORTENER_*` environment variables.
- [ ] **Create `Caddyfile`**:
    - Configure reverse proxy for the application.
    - Use a placeholder domain (e.g., `url-shortener.sslip.io` or just `localhost` for now) until OCI IP is known.

### 3. Build & Run Scripts
- [ ] **Update `Dockerfile`**:
    - Change base image to `gcr.io/distroless/static-debian12:nonroot`.
    - Expose port `8080`.
    - Ensure `USER` is set to non-root.
- [ ] **Update `Makefile`**:
    - Replace manual docker commands with `docker compose` commands (`up`, `down`, `logs`).
    - Simplify `run` and `db` targets.

### 4. CI/CD
- [ ] **Update `.github/workflows/aws.yml`**:
    - Rename to `.github/workflows/deploy.yml`.
    - Update "Deploy" job to use OCI/Podman logic from `go-dutch`.
    - Steps:
        - Setup SSH.
        - Configure Cloud Instance (Firewall, Podman socket, Sysctl).
        - Clone/Pull repository.
        - Create `.env` file dynamically.
        - Run `docker compose build`, `docker compose up migrate`, `docker compose up -d`.
        - Verify deployment.

### 5. Documentation
- [ ] **Update `README.md`**:
    - Update deployment instructions for OCI.
    - Update local development instructions to use Docker Compose.
    - Remove AWS-specific mentions and Nginx configuration details.

## Next Steps
- Awaiting user confirmation to execute the file modifications.