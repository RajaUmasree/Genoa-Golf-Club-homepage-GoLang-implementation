# GGC_GL_01-3 — Root Cause & Backend Analysis: Homepage

**Project:** Genoa Golf Club (GGC)

## 1. Problem Statement

The marketing site currently has no dedicated backend for homepage content —
marketing staff cannot update the hero banner or announcements without a code
change and redeploy. This ticket investigates the backend gap and proposes the
fix implemented in 01-5.

## 2. Root Cause

- No `homepage` domain existed in the GGC Go backend; hero/announcement text
  was hardcoded in the front-end bundle.
- Events and membership tiers already exist as domains elsewhere in the GGC
  backend, but there was no aggregation layer to combine them with
  homepage-only content (hero, announcements, gallery) into one payload for
  the landing page, forcing (hypothetically) multiple front-end round trips.

## 3. Impact

| Area | Impact |
|---|---|
| Marketing agility | Every banner/announcement change required a dev + deploy cycle |
| Front-end performance | Multiple separate calls needed to render one page |
| Consistency | No single source of truth / audit trail for homepage content changes |

## 4. Backend Design Decision

- Introduce a new `homepage` bounded context inside the existing Go backend
  (matches current architecture: Gin + MongoDB + JWT, per `genoa-golf-club`
  backend conventions).
- Repository pattern with an interface (`HomepageRepository`) so the storage
  implementation is swappable — this deliverable ships:
  - `InMemoryHomepageRepository` (default, zero external dependencies, used
    for local dev / this ticket's demo & tests)
  - A `mongo` build tag stub showing the intended MongoDB-backed
    implementation, consistent with the rest of the GGC backend, to be wired
    to the live cluster during deployment (01-8).
- Aggregation endpoint (`GET /api/v1/homepage`) implemented in the service
  layer via fan-out to the same repository, avoiding duplicate network calls
  once this is deployed alongside the existing Events/Membership services.

## 5. Alternatives Considered

| Option | Rejected because |
|---|---|
| Let front-end call 5 separate services directly | More round trips, more front-end complexity, harder to cache |
| Store homepage content in a CMS (e.g. Contentful) | Out of scope/budget for this epic; revisit later |
| Hardcode content and redeploy for changes (status quo) | Root cause of this ticket — rejected |

## 6. Risks & Mitigations

| Risk | Mitigation |
|---|---|
| In-memory repo loses data on restart | Documented limitation; Mongo implementation is the production path (01-8 runbook) |
| Admin routes could be abused if JWT secret leaks | Reuses existing GGC JWT middleware/secret rotation process |
| Aggregation endpoint becomes a hot path | Cheap in-process fan-out; safe to add response caching later without API changes |
