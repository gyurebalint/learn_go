# Crypto Aggregator API

A containerized Go service that fetches cryptocurrency prices and persists them to PostgreSQL.

## Core Functionality
- **API:** REST endpoint for crypto price retrieval.
- **Resilience:** Implements database connection retries with exponential backoff.
- **Config:** Environment-variable driven configuration.
- **Build:** Multi-stage Dockerfile for optimized image size.

## Technical Stack
- **Language:** Go 1.25
- **Database:** PostgreSQL 15
- **Infrastructure:** Docker, Docker Compose

## Quick Start (Docker Compose)
The easiest way to run the full stack (App + Database) is using Docker Compose:

```bash
docker-compose up --build
```
This will initialize the PostgreSQL database and start the aggregator service on `localhost:8080`.

[Image of Docker Compose architecture for microservices]

## Environment Variables
- `DB_HOST`: Database hostname (default: `localhost`)
- `DB_USER`: Database user (default: `postgres`)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `PORT`: Service port (default: `3000`)

## Manual Usage (Docker Only)
If you already have a Postgres instance running:

```bash
docker build -t crypto-aggregator:v2 .
docker run -p 8080:3000 --env DB_HOST=host.docker.internal crypto-aggregator:v2
```

## API Endpoint
`GET /price?coin=bitcoin`

**Response:**
```json
{"coin": "bitcoin", "price": "98000.50", "currency": "USD"}
```