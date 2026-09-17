# GGC_GL_01-1 — Interface & Requirement Analysis: Homepage

**Project:** Genoa Golf Club (GGC)
**Module:** Homepage
**Owner:** Umasree
**Status:** Complete

## 1. Purpose

Define the functional and interface requirements for the Genoa Golf Club public-facing
homepage backend, which powers the marketing site's landing page (hero banner,
announcements, upcoming events, membership highlights, and photo gallery).

## 2. Stakeholders

| Role | Interest |
|---|---|
| Club Marketing Team | Controls hero banner, announcements, promotions |
| Membership Office | Publishes membership tier highlights |
| Events/Pro Shop | Publishes upcoming tournaments and events |
| Front-end team | Consumes the Homepage API to render the site |
| Club Members / Public Visitors | End consumers of the homepage content |

## 3. Business Requirements

- BR-1: Visitors must see an always-current hero banner (image, headline, CTA).
- BR-2: Visitors must see the 3 most recent club announcements.
- BR-3: Visitors must see up to 5 upcoming events, soonest first.
- BR-4: Visitors must see the active membership tiers with a summary of benefits.
- BR-5: Visitors must see a curated photo gallery.
- BR-6: Marketing/admin staff must be able to update hero banner and announcement
  content without a code deployment (via authenticated API).
- BR-7: The homepage must load as a single aggregated payload to minimize
  front-end round trips.

## 4. Interface Requirements

### 4.1 Consumer-facing (public, read-only)
- `GET /api/v1/homepage` — aggregated homepage payload (single call for the site).
- `GET /api/v1/homepage/hero`
- `GET /api/v1/homepage/announcements`
- `GET /api/v1/homepage/events`
- `GET /api/v1/homepage/membership-highlights`
- `GET /api/v1/homepage/gallery`

### 4.2 Admin-facing (JWT-protected, CMS-style management)
- `POST /api/v1/admin/homepage/hero`
- `PUT  /api/v1/admin/homepage/hero/:id`
- `POST /api/v1/admin/homepage/announcements`
- `DELETE /api/v1/admin/homepage/announcements/:id`

## 5. Non-Functional Requirements

- NFR-1: p95 latency < 150ms for the aggregated homepage endpoint (cache-friendly).
- NFR-2: All admin write endpoints require a valid JWT bearer token (existing
  GGC auth service, HS256).
- NFR-3: Service must be horizontally stateless (state lives in MongoDB), so it
  fits the existing GGC Go backend's containerized deployment model.
- NFR-4: All timestamps in UTC / RFC3339.
- NFR-5: Structured JSON logging; no PII in logs.

## 6. Out of Scope (this Jira epic)

- Front-end rendering of the homepage (handled by web team).
- Payment/booking flows (covered by separate tee-time and membership epics).
- Image upload/transcoding pipeline (assumed pre-existing CDN asset URLs are supplied).

## 7. Acceptance Criteria

- [ ] All 5 read endpoints return 200 with valid JSON matching the API contract (01-4).
- [ ] Aggregated `/homepage` endpoint returns all sections in one response.
- [ ] Admin write endpoints reject requests without a valid JWT (401).
- [ ] Unit + integration tests pass (see 01-7).
