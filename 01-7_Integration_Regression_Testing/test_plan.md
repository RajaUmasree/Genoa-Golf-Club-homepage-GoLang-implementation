# GGC_GL_01-7 — Integration / Regression Testing: Homepage

**Project:** Genoa Golf Club (GGC)

## 1. Test Levels

| Level | Location | Tooling |
|---|---|---|
| Unit — repository | `internal/repository/homepage_repository_test.go` | `testing` |
| Unit — service (business rules) | `internal/service/homepage_service_test.go` | `testing` |
| Integration — full HTTP stack (routes + middleware + handlers) | `internal/handlers/homepage_handler_test.go` | `net/http/httptest` |

The integration tests boot the real `http.ServeMux` produced by
`routes.RegisterRoutes`, so they exercise JWT middleware, routing, and JSON
(de)serialization exactly as they run in production — not mocks.

## 2. Test Matrix

| # | Scenario | Type | Result |
|---|---|---|---|
| 1 | Seeded hero banner is retrievable | Unit | PASS |
| 2 | Upserting a hero banner updates the active hero | Unit | PASS |
| 3 | Announcements are limited and sorted newest-first | Unit | PASS |
| 4 | Creating then deleting an announcement; second delete → `ErrNotFound` | Unit | PASS |
| 5 | Past events are excluded from "upcoming events" | Unit | PASS |
| 6 | Upcoming events are sorted ascending and respect `limit` | Unit | PASS |
| 7 | Membership highlights sorted ascending by price | Unit | PASS |
| 8 | Gallery images are returned | Unit | PASS |
| 9 | Aggregated homepage payload populates all 5 sections | Service | PASS |
| 10 | `ListAnnouncements(0)` falls back to default limit (3) | Service | PASS |
| 11 | `ListUpcomingEvents(0)` falls back to default limit (5) | Service | PASS |
| 12 | Setting a new hero marks it active and persists | Service | PASS |
| 13 | Updating a hero with a mismatched ID returns `ErrNotFound` | Service | PASS |
| 14 | Create/delete announcement round-trip via service | Service | PASS |
| 15 | `GET /healthz` → 200 | Integration | PASS |
| 16 | `GET /api/v1/homepage` → 200, all sections present | Integration | PASS |
| 17 | `GET /api/v1/homepage/announcements?limit=1` → exactly 1 item | Integration | PASS |
| 18 | `POST /admin/homepage/hero` without token → 401 | Integration | PASS |
| 19 | `POST /admin/homepage/hero` with expired token → 401 | Integration | PASS |
| 20 | `POST /admin/homepage/hero` with valid token → 201, and public `GET /hero` reflects the change | Integration | PASS |
| 21 | `DELETE /admin/homepage/announcements/{unknown-id}` with valid token → 404 | Integration | PASS |
| 22 | `GET /api/v1/homepage/events` never returns past events | Integration | PASS |

**Total: 22 test cases, 22 passing (0 failing).**

## 3. Regression Coverage

Re-running the full suite is the regression gate for any future change to
this module:

```bash
cd 01-5_GoLang_Implementation/homepage-service
go test ./... -v -cover
```

Any PR touching `internal/homepage_service` domain code must keep all 22
cases green before merge (enforced in CI — see 01-8 for the pipeline
description).

## 4. Manual / Exploratory Testing Performed

- Ran the compiled binary locally (`go build -o /tmp/homepage-service
  ./cmd/server`) and hit it with `curl`:
  - `GET /healthz` → `{"status":"ok"}`
  - `GET /api/v1/homepage` → full aggregated JSON payload
  - `POST /api/v1/admin/homepage/hero` without a token → `401`
- Confirmed no goroutine leaks or panics under the above manual exercise.

## 5. Known Gaps / Follow-ups

- No load/performance test yet against the NFR-1 target (p95 < 150ms) —
  recommend a follow-up ticket once deployed behind the real ingress, since
  in-process latency in this deliverable is sub-millisecond and not
  representative of the network path.
- Mongo-backed repository (`-tags mongo`) has no automated tests yet since it
  is a documented stub pending driver wiring (see 01-8).
