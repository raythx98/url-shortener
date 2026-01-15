# url-shortener

## Tech Stack

Reverse Proxy: Caddy (Automatic HTTPS)

Backend: Go, Postgres

CA: Let's Encrypt (via Caddy)

## Design Considerations

### Architecture
The architecture has been migrated to a containerized setup orchestrated by Docker Compose, running on Oracle Cloud Infrastructure (OCI).

- **Caddy**: Acts as the reverse proxy and handles automatic HTTPS certificate management.
- **App**: The Go backend service.
- **DB**: PostgreSQL database.
- **Migrate**: Runs database migrations on startup.

This setup is deployed on a single OCI instance using Podman (daemonless Docker alternative) for enhanced security and `docker-compose` for orchestration.

### Short URL Collision

An alphanumeric character set typically includes:
- 26 uppercase letters (A-Z)
- 26 lowercase letters (a-z)
- 10 digits (0-9)

Total: 62 characters

The current url shortener is configured to use a 10-character alphanumerical string, 
the total number of possible unique values is:

$N = 62^{10}$

$N \approx 8.39 \times 10^{17}$

Using the Birthday Paradox Approximation, the probability of at least one collision among
𝑘 randomly chosen values is:

$P \approx 1 - e^{-\frac{k^2}{2N}}$

#### For 1 million randomly generated URLs, the probability of a collision is negligible

$P \approx 1 - e^{-\frac{(10^6)^2}{2 \times 8.39 \times 10^{17}}}$

$P \approx 0.00006$%

#### For 1 billion randomly generated URLs, the probability of a collision is significant
$P \approx 1 - e^{-\frac{(10^9)^2}{2 \times 8.39 \times 10^{17}}}$

$P \approx 44.9$%

To mitigate this risk when site traffic increases, we can:
- Increase ID length (e.g., use 12+ characters).
- Implement a database uniqueness check to prevent collisions.

## Infrastructure Setup 

### GitHub Secrets
You will need to set up the following in Github Secrets:
- `BASIC_AUTH_PASSWORD`
- `DB_PASSWORD`
- `DEPLOY_HOST` (OCI Instance IP)
- `DEPLOY_USER` (OCI Instance Username, e.g., `opc`)
- `JWT_SECRET`
- `SSH_PRIVATE_KEY`

### Generate Secrets
An easy way to generate cryptographically secure random strings is to use the following command:
```bash
python -c 'import secrets; print(secrets.token_urlsafe(32))'
```

### OCI Instance Setup (Automated via Workflow)
The GitHub Action workflow automatically handles most of the instance configuration, including:
- Updating system packages.
- Configuring the firewall (ports 80, 443).
- Installing and configuring Podman to emulate Docker.
- Enabling rootless containers to bind to privileged ports.
- Installing the Docker Compose plugin.

### Local Development
To run the project locally:

1.  **Environment Variables**: Create a `.env` file based on the example in the CI/CD workflow or `Makefile` defaults (though `docker-compose.yaml` uses `.env`).

2.  **Run with Docker Compose**:
    ```bash
    make up
    ```
    This will start the App, DB, Migration, and Caddy services.

3.  **View Logs**:
    ```bash
    make logs
    ```

4.  **Stop**:
    ```bash
    make down
    ```

### Manual Deployment (Reference)
The deployment is handled by `.github/workflows/deploy.yml`. It uses SSH to connect to the OCI instance, pulls the latest code, builds the images using Podman (via Docker alias), and restarts the services using `docker compose`.

```
