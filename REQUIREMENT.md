# Requirements

## Active Task
- [x] Migrate deployment from AWS to OCI.
    - [x] Remove Nginx configuration and dependency.
    - [x] Modify `aws.yml` (rename to `deploy.yml`) to deploy to OCI, to match `go-dutch` reference.
    - [x] Update `Dockerfile` to match `go-dutch` reference.
    - [x] Update `Makefile` to match `go-dutch` reference.
    - [x] Update `docker-compose.yaml` (or equivalent) to match `go-dutch` reference.
    - [x] Update `README.md` to reflect OCI migration and new deployment/run instructions.
    - [x] Update Go version to 1.24.11 across the project (`go.mod`, `Dockerfile`, CI/CD).
- [x] Rename repository from `url-shortener` to `go-zap`.
    - [x] Update `go.mod` module path.
    - [x] Update all internal imports.
    - [x] Restructure mocks directory.

## Constraints
- Align configuration with `../go-dutch` project.