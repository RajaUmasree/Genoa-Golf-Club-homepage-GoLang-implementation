# GGC_GL_01-4 — Technical Design / API Contract: Homepage

**Project:** Genoa Golf Club (GGC)

## 1. Package Layout

```
homepage-service/
├── cmd/server/main.go              # entrypoint, wiring
├── internal/
│   ├── config/config.go            # env config loader
│   ├── models/homepage.go          # domain structs
│   ├── repository/
│   │   ├── homepage_repository.go  # interface + in-memory impl
│   │   └── homepage_repository_test.go
│   ├── service/
│   │   ├── homepage_service.go     # business rules
│   │   └── homepage_service_test.go
│   ├── handlers/
│   │   ├── homepage_handler.go     # Gin handlers
│   │   └── homepage_handler_test.go
│   ├── middleware/auth.go          # JWT auth middleware
│   └── routes/routes.go            # route registration
├── go.mod
├── Dockerfile
├── Makefile
└── README.md
```

## 2. Domain Models

```go
type HeroBanner struct {
    ID        string    `json:"id" bson:"_id"`
    Headline  string    `json:"headline" bson:"headline"`
    Subtext   string    `json:"subtext" bson:"subtext"`
    ImageURL  string    `json:"imageUrl" bson:"imageUrl"`
    CTALabel  string    `json:"ctaLabel" bson:"ctaLabel"`
    CTAUrl    string    `json:"ctaUrl" bson:"ctaUrl"`
    Active    bool      `json:"active" bson:"active"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Announcement struct {
    ID        string    `json:"id" bson:"_id"`
    Title     string    `json:"title" bson:"title"`
    Body      string    `json:"body" bson:"body"`
    PostedAt  time.Time `json:"postedAt" bson:"postedAt"`
}

type Event struct {
    ID        string    `json:"id" bson:"_id"`
    Name      string    `json:"name" bson:"name"`
    StartDate time.Time `json:"startDate" bson:"startDate"`
    Location  string    `json:"location" bson:"location"`
}

type MembershipHighlight struct {
    ID       string   `json:"id" bson:"_id"`
    Tier     string   `json:"tier" bson:"tier"`
    Price    float64  `json:"price" bson:"price"`
    Benefits []string `json:"benefits" bson:"benefits"`
}

type GalleryImage struct {
    ID      string `json:"id" bson:"_id"`
    URL     string `json:"url" bson:"url"`
    Caption string `json:"caption" bson:"caption"`
}

type HomepageResponse struct {
    Hero        *HeroBanner           `json:"hero"`
    Announcements []Announcement      `json:"announcements"`
    UpcomingEvents []Event            `json:"upcomingEvents"`
    Membership  []MembershipHighlight `json:"membershipHighlights"`
    Gallery     []GalleryImage        `json:"gallery"`
}
```

## 3. API Contract (OpenAPI-style summary)

### GET /api/v1/homepage
200 → `HomepageResponse`

### GET /api/v1/homepage/hero
200 → `HeroBanner` | 404

### GET /api/v1/homepage/announcements?limit=3
200 → `Announcement[]`

### GET /api/v1/homepage/events?limit=5
200 → `Event[]` (only future events, ascending by startDate)

### GET /api/v1/homepage/membership-highlights
200 → `MembershipHighlight[]`

### GET /api/v1/homepage/gallery
200 → `GalleryImage[]`

### POST /api/v1/admin/homepage/hero  (JWT required)
Body: `{headline, subtext, imageUrl, ctaLabel, ctaUrl}` → 201 → `HeroBanner`

### PUT /api/v1/admin/homepage/hero/:id  (JWT required)
Body: partial `HeroBanner` fields → 200 → `HeroBanner` | 404

### POST /api/v1/admin/homepage/announcements  (JWT required)
Body: `{title, body}` → 201 → `Announcement`

### DELETE /api/v1/admin/homepage/announcements/:id  (JWT required)
204 | 404

### GET /healthz
200 → `{"status":"ok"}`

## 4. JWT Middleware

- Reads `Authorization: Bearer <token>` header.
- Validates HS256 signature against `JWT_SECRET` env var.
- On failure → 401 with the standard error envelope.
- Full OpenAPI YAML: see `openapi.yaml` in this same folder.
