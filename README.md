# mood-api

A microservices-based mood tracking and advice API built with Go. The system provides mood tracking, personalized advice recommendations, quote management, and authentication through a unified API gateway.

## Overview

`mood-api` is a distributed system that separates concerns into independent microservices. The **API Gateway** acts as the single entry point for clients, routing requests to backend services and handling cross-service orchestration.

### Architecture

- **API Gateway** (`gateway`): Routes requests, authenticates users, and orchestrates multi-service workflows
- **Mood Service** (`mood`): Manages mood entries and mood summaries
- **Advice Service** (`advice`): Selects and tracks advice recommendations based on mood patterns
- **Auth Service** (`auth`): User registration, login, token management
- **Quote Service** (`quote`): Fetches and caches daily motivational quotes (Quotes provided by: https://zenquotes.io/)
- **PostgreSQL**: Shared database for mood, advice, auth, and user data
- **Redis**: Cache layer for quote data

## Usage

### Requirements

- Docker & Docker Compose installed ([download here](https://www.docker.com/products/docker-desktop))

### Quick Start

#### Option 1: With Git (Recommended)

```bash
# Clone the repository
git clone https://github.com/ciameksw/mood-api.git
cd mood-api

# Start all services (builds images locally)
docker compose up -d --build
```

#### Option 2: Without Git

1. Download the repository as ZIP from GitHub
2. Extract the folder
3. Open terminal in the extracted folder
4. Run:
```bash
docker compose up -d --build
```

### Stop the System

```bash
docker compose down
```

This removes all containers but keeps data in PostgreSQL.

### Clean Up Everything

```bash
docker compose down -v
```

This removes containers and volumes (deletes all data).

### Gateway API

All requests go through the gateway.

See [GATEWAY_API.md](./GATEWAY_API.md) for detailed endpoint documentation.

## Ownership

Built and maintained by @ciameksw.

## License

Distributed under the MIT License. See [LICENSE](./LICENSE) for details.