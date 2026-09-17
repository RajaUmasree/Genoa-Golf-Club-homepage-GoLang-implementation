# Genoa Golf Club — Homepage Service

Backend service for the GGC marketing site homepage (Jira epic `GGC_GL_01`).
Serves the hero banner, announcements, upcoming events, membership
highlights, and gallery for the public landing page, plus JWT-protected
admin endpoints for content management.

Full SDLC documentation for each Jira sub-task lives one level up in the
repo (`../../01-1_Interface_Requirement_Analysis` etc.) — see the top-level
`GGC_GL_01_Homepage/README.md` for the index.

## Requirements

- Go 1.22+

## Run locally

```bash
go run ./cmd/server
# service listens on :8080 by default
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/homepage
```

## Configuration (environment variables)

| Var | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `JWT_SECRET` | `dev-only-secret-change-me` | HS256 secret for admin routes — **override in every non-local environment** |
| `MONGO_URI` | `mongodb://localhost:27017` | Used only when built with `-tags mongo` |
| `MONGO_DB` | `genoa_golf_club` | Used only when built with `-tags mongo` |

## Storage

Ships with a thread-safe **in-memory repository** (seeded with representative
data) so it runs with zero external dependencies. A MongoDB-backed
implementation matching the existing GGC backend's conventions is stubbed in
`internal/repository/homepage_repository_mongo.go` (build tag `mongo`) — swap
it in for production once the collections are provisioned (see
`../01-8_Evidence_Defect_Closure_Release_Readiness/release_readiness.md`).

## Test

```bash
make test
# or
go test ./... -v -cover
```

## Build & run with Docker

```bash
make docker-build
docker run -p 8080:8080 -e JWT_SECRET=change-me genoagolfclub/homepage-service:latest
```

## API

See `../01-4_Technical_Design_API_Contract/openapi.yaml` for the full contract.
