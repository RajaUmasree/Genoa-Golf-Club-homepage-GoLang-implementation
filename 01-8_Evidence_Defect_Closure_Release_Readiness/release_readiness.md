# GGC_GL_01-8 — Evidence / Defect Closure / Release Readiness: Homepage

**Project:** Genoa Golf Club (GGC)

## 1. Evidence

See `test_evidence.txt` in this folder for the raw `go build` / `go vet` /
`go test -v -cover` output captured at sign-off time. Summary:

| Check | Result |
|---|---|
| `go build ./...` | PASS (no output, zero errors) |
| `go vet ./...` | PASS (no issues) |
| `go test ./...` | PASS — 22/22 test cases green |
| Coverage | repository 95.0%, service 67.3%, handlers 56.1% |
| Manual smoke test (`curl` against running binary) | PASS — see 01-7 §4 |

## 2. Defect Log

| ID | Description | Severity | Status |
|---|---|---|---|
| DEF-1 | Initial design assumed Gin + golang-jwt (per existing GGC backend convention), but the build environment used for this deliverable could not reach the Go module proxy | Medium | **Closed** — reimplemented on Go 1.22 stdlib (`net/http` + hand-rolled HS256 JWT check), fully offline-buildable, functionally equivalent. See 01-6 §3 for the review note. |
| DEF-2 | Import cycle in initial handler integration test (`handlers` test importing `routes`, which imports `handlers`) | Low | **Closed** — moved integration tests to external test package `handlers_test`. |

No open/unresolved defects at time of sign-off.

## 3. Release Readiness Checklist

- [x] Code builds cleanly (`go build ./...`)
- [x] Static analysis clean (`go vet ./...`)
- [x] Full automated test suite passing (22/22)
- [x] JWT-protected admin routes verified to reject missing/expired tokens
- [x] Dockerfile present and uses a minimal, non-root runtime image
- [x] `.gitignore` present (excludes `bin/`, `.env`, IDE folders)
- [x] README documents run/build/test/Docker instructions and all env vars
- [x] OpenAPI contract (01-4) matches implemented routes
- [ ] MongoDB-backed repository wired and load-tested against staging cluster
      — **tracked as a follow-up**, not a blocker for this ticket's scope
      (in-memory repository is a documented, intentional interim state; see
      01-3 §4 and README "Storage" section)
- [ ] CI pipeline definition (GitHub Actions) — **recommended next step**,
      see §4 below for a ready-to-use starter workflow

## 4. Suggested CI Workflow (for the target GitHub repo)

Add this as `.github/workflows/ci.yml` once pushed to GitHub — it runs the
exact same checks captured as evidence above on every push/PR:

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:

jobs:
  build-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: 01-5_GoLang_Implementation/homepage-service
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./... -v -cover
```

(Adjust `working-directory` if the Go module is moved/flattened in the
target repository layout.)

## 5. Deployment Runbook (summary)

1. `docker build -t genoagolfclub/homepage-service:<tag> .` from
   `01-5_GoLang_Implementation/homepage-service`.
2. Push to the club's container registry.
3. Set `JWT_SECRET` (shared with the existing GGC auth service) and, once the
   Mongo implementation is enabled, `MONGO_URI` / `MONGO_DB` as deployment
   secrets/env vars — **never** rely on the in-code dev-only default.
4. Roll out behind the existing ingress alongside the other GGC backend
   services; `/healthz` is the liveness/readiness probe path.
5. Point the marketing site's front-end at `GET /api/v1/homepage` for the
   landing page's single aggregated data call.

## 6. Sign-off

**Status: Ready for release**, with the two open follow-ups in §3 tracked as
separate, non-blocking work items (Mongo cutover, CI workflow activation).
