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
Run the full stack (App + Database) with one command:

```bash
docker-compose up --build
```



## Verifying Data Storage
Once the application is running and you have made a request to `/price`, you can verify the data is saved in PostgreSQL using the following methods:

### 1. Via Docker (If using Docker Compose)
Access the database directly inside the running container:
```bash
# Find the container name
docker ps

# Enter the database
docker exec -it <postgres_container_name> psql -U admin -d crypto-aggregator

# Run the query
SELECT * FROM prices;
```

### 2. Via Local Terminal (If running app locally)
If you are running the app outside of Docker, use a local SQL client (psql, DBeaver, or pgAdmin) to connect to `localhost:5432`.



## Environment Variables
- `DB_HOST`: Database hostname (default: `localhost`)
- `DB_USER`: Database user (default: `postgres`)

## API Endpoint
`GET /price?symbol=bitcoin`

**Response:**
```json
{"coin": "bitcoin", "price": "98000.50", "currency": "USD"}
```