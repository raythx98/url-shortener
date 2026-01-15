# Requirements

## Active Task
- [ ] Migrate deployment from AWS to OCI.
    - [ ] Remove Nginx configuration and dependency.
    - [ ] Modify `aws.yml` (rename to `deploy.yml`) to deploy to OCI, to match `go-dutch` referencce.
    - [ ] Update `Dockerfile` to match `go-dutch` reference.
    - [ ] Update `Makefile` to match `go-dutch` reference.
    - [ ] Update `docker-compose.yaml` (or equivalent) to match `go-dutch` reference.
    - [ ] Update `README.md` to reflect OCI migration and new deployment/run instructions.
    - [ ] Update Go version to 1.24.11 across the project (`go.mod`, `Dockerfile`, CI/CD).

## Constraints
- Align configuration with `../go-dutch` project.
- Do not make code changes yet (Phase 1: Plan & Document).
