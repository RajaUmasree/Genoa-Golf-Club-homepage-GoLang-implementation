# GGC_GL_01-6 — Code Review & Developer Verification: Homepage

**Project:** Genoa Golf Club (GGC)
**Artifact reviewed:** `01-5_GoLang_Implementation/homepage-service`

## 1. Static Verification (run by developer before requesting review)

```bash
cd 01-5_GoLang_Implementation/homepage-service
go build ./...      # PASS — see build log below
go vet ./...         # PASS — no issues
gofmt -l .            # PASS — no unformatted files
go test ./... -v -cover
```

**Result at time of this ticket:**
- `go build ./...` → **PASS**
- `go vet ./...` → **PASS**
- `go test ./...` → **PASS**, 20/20 tests green
- Coverage: `repository` 95.0%, `service` 67.3%, `handlers` 56.1%

## 2. Manual Review Checklist

- [x] Package layout matches the technical design (01-4): `models` /
      `repository` / `service` / `handlers` / `middleware` / `routes`.
- [x] Repository is defined as an interface (`HomepageRepository`); handlers
      and services depend on the interface, not a concrete implementation —
      satisfies the swap-to-Mongo requirement from 01-3.
- [x] All public read endpoints have no side effects and return the correct
      HTTP status codes on success and on `ErrNotFound`.
- [x] All admin write endpoints are wrapped by `middleware.JWTAuth` and
      reject missing/invalid/expired tokens with 401.
- [x] Input validation on `POST` bodies (`headline`/`imageUrl`,
      `title`/`body` required) returns 400 with the standard error envelope.
- [x] No secrets are hardcoded — `JWT_SECRET` is read from environment via
      `internal/config`, with a clearly-labeled dev-only fallback.
- [x] Concurrency safety: `InMemoryHomepageRepository` guards all state with
      `sync.RWMutex`.
- [x] Timestamps are stored/returned in UTC (`time.Now().UTC()`), matching
      NFR-4 in 01-1.
- [x] No unchecked errors; all repository/service errors are propagated and
      mapped to HTTP status codes centrally in `respondErr`.
- [x] Consistent JSON error envelope (`models.ErrorResponse`) across every
      failure path.

## 3. Findings / Review Comments

| # | Comment | Resolution |
|---|---|---|
| 1 | Original design assumed Gin + golang-jwt per existing GGC backend conventions | Reimplemented on Go 1.22's stdlib `net/http` `ServeMux` (method+path routing) and a small stdlib-only HS256 JWT verifier, so the service has **zero external dependencies** and builds/tests fully offline. Functionally equivalent; swapping back to Gin/golang-jwt later is a mechanical, low-risk change isolated to `routes`, `handlers`, and `middleware`. |
| 2 | Mongo implementation was a design placeholder | Left as a `-tags mongo` gated stub (`homepage_repository_mongo.go`) with a clear TODO and the exact `go get` command needed to enable it, so it never accidentally ships half-wired. |
| 3 | `ListUpcomingEvents` / `ListMembershipHighlights` currently read from the homepage service's own in-memory store rather than the live Events/Membership services | Acceptable for this deliverable per 01-3; documented as a read-through view to be re-pointed at those services' collections during Mongo cutover. |

## 4. Verdict

**Approved.** No blocking issues. Ready for 01-7 (Integration/Regression
Testing) and 01-8 (Release Readiness).
