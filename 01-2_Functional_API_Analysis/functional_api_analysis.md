# GGC_GL_01-2 — Functional & API Analysis: Homepage

**Project:** Genoa Golf Club (GGC)
**Depends on:** GGC_GL_01-1 (Requirements)

## 1. Functional Breakdown

| Function | Description | Priority |
|---|---|---|
| F-1 Hero Banner Retrieval | Return the currently active hero banner | Must |
| F-2 Announcements Feed | Return latest N announcements, newest first | Must |
| F-3 Upcoming Events Feed | Return events with `startDate >= now`, sorted ascending | Must |
| F-4 Membership Highlights | Return active membership tiers with benefit summaries | Must |
| F-5 Gallery | Return curated gallery images with captions | Should |
| F-6 Aggregation | Combine F-1..F-5 into a single homepage payload | Must |
| F-7 Admin CRUD (hero) | Create/update hero banner content | Must |
| F-8 Admin CRUD (announcements) | Create/delete announcements | Should |

## 2. API Surface (summary — full contract in 01-4)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| GET | /api/v1/homepage | none | Aggregated payload (F-6) |
| GET | /api/v1/homepage/hero | none | F-1 |
| GET | /api/v1/homepage/announcements | none | F-2 |
| GET | /api/v1/homepage/events | none | F-3 |
| GET | /api/v1/homepage/membership-highlights | none | F-4 |
| GET | /api/v1/homepage/gallery | none | F-5 |
| POST | /api/v1/admin/homepage/hero | JWT | F-7 create |
| PUT | /api/v1/admin/homepage/hero/:id | JWT | F-7 update |
| POST | /api/v1/admin/homepage/announcements | JWT | F-8 create |
| DELETE | /api/v1/admin/homepage/announcements/:id | JWT | F-8 delete |
| GET | /healthz | none | liveness probe |

## 3. Data Ownership

- Homepage service owns its own collections (`hero_banners`, `announcements`,
  `gallery_images`), each pre-fixed `homepage_` in Mongo to avoid collisions
  with other GGC modules (membership, events) already in the shared cluster.
- `events` and `membership-highlights` are treated as **read-through** views:
  in this deliverable they use the same repository abstraction so they can
  later be pointed at the existing Events/Membership services' collections
  instead of duplicating data — see 01-3 (Root Cause & Backend Analysis) for
  the rationale.

## 4. Error Model

Standard envelope used across all endpoints:

```json
{ "error": { "code": "NOT_FOUND", "message": "hero banner not found" } }
```

| HTTP | Code | When |
|---|---|---|
| 400 | VALIDATION_ERROR | malformed request body |
| 401 | UNAUTHORIZED | missing/invalid JWT on admin routes |
| 404 | NOT_FOUND | resource id not found |
| 500 | INTERNAL_ERROR | unexpected server/repo error |

## 5. Sequencing / Dependencies

1. Models → 2. Repository interfaces (+ in-memory impl for this deliverable,
Mongo impl to follow existing GGC pattern) → 3. Service layer (business rules:
filtering upcoming events, limiting announcement count) → 4. Handlers → 5. Routes
→ 6. Middleware (JWT) wiring → 7. Tests.
