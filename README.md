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

### Start the System

```bash
docker compose up -d --build
```

Gateway service will be available at: `http://localhost:3000`

### Gateway API

All requests go through the gateway.

See [GATEWAY_API.md](./GATEWAY_API.md) for detailed endpoint documentation.