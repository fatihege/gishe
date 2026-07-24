# Gishe

A ticketing backend focused on correctness under concurrent demand.

> Status: Work in progress. Basic authentication and catalog management are implemented.
> Writing tests and solidifying the current implementations are the next milestones.

## Why this project exists?

I'm building this project to understand the concept of correctness under concurrent demand. So it is a learning-focused project about transactional seat reservations, idempotent payments, and reliable event processing.

I utilize AI as little as possible and implement especially the core mechanisms myself in order to build a solid understanding of the architectural design and implementation of such systems.

## Current features

- PostgreSQL migrations
- User registration and login
- JWT authentication
- Venue creation and listing
- Session creation and listing

## Planned milestones

- Writing tests and solidifying the implementation so far
- Seat creation, listing, and seeding
- PostgreSQL-based seat holds
- Concurrent hold verification
- Orders and mock payments
- Idempotency and ledger
- Outbox and asynchronous ticket creation

## Architecture

```
Client → HTTP Handler → Application Service → Repository → PostgreSQL
```

## Tech stack

- Go
- chi
- PostgreSQL
- pgx
- goose

## Getting started

Follow the steps below in order to run the project on your machine.

### Prerequisites

- Git
- Go 1.26.5+
- Docker
- goose 3.27.2+

### Installation

1. Clone the repository and navigate into it.
```bash
git clone https://github.com/fatihege/gishe
cd gishe
```

2. You first need to copy the `.env.example` file as `.env` and assign values as you desired. You can see the environment configuration below.
```bash
cp .env.example .env
```

3. And download Go dependencies.
```bash
go mod download
```

### Usage

1. First of all, you need to start Docker containers by providing the `.env` file and the compose file.
```bash
docker compose --env-file .env -f deploy/docker-compose.yml up -d
```

2. Make sure to source and load the environment variables before running the application. One way to do it:
```bash
set -a
source .env
set +a
```

3. Don't forget to apply migrations. If you set goose environment variables (see Configuration), you can simply run:
```bash
goose up
```

4. Finally, start the Go application.
```bash
go run ./cmd/api
```

## Configuration

| Variable | Required | Default | Purpose |
|---|---:|---|---|
| POSTGRES_DB | Yes | - | PostgreSQL database name |
| POSTGRES_USER | Yes | - | PostgreSQL user |
| POSTGRES_PASSWORD | No | - | PostgreSQL password |
| POSTGRES_HOST | Yes | - | PostgreSQL host |
| POSTGRES_PORT | Yes | - | PostgreSQL port |
| POSTGRES_SSL_MODE | No | disable | PostgreSQL connection SSL mode |
| GOOSE_DRIVER | No | - | goose DB target |
| GOOSE_DBSTRING | No | - | goose connection string |
| GOOSE_MIGRATION_DIR | No | - | goose migrations folder |
| HTTP_ADDRESS | No | :8080 | Server bind address |
| JWT_SECRET | Yes | - | JWT signing secret |

## Running tests

Tests are not created yet, as you can see.
