// External test package (handlers_test) — this test exercises the full
// route+middleware+handler stack via httptest, which would otherwise create
// an import cycle (routes imports handlers) if declared as part of the
// internal "handlers" package.
package handlers_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/genoagolfclub/homepage-service/internal/handlers"
	"github.com/genoagolfclub/homepage-service/internal/models"
	"github.com/genoagolfclub/homepage-service/internal/repository"
	"github.com/genoagolfclub/homepage-service/internal/routes"
	"github.com/genoagolfclub/homepage-service/internal/service"
)

const testJWTSecret = "test-secret"

func newTestMux() *http.ServeMux {
	repo := repository.NewInMemoryHomepageRepository()
	svc := service.NewHomepageService(repo)
	h := handlers.NewHomepageHandler(svc)
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, h, testJWTSecret)
	return mux
}

// makeTestJWT builds a minimal valid HS256 JWT for the given secret so admin
// routes can be exercised end-to-end in tests without an external library.
func makeTestJWT(secret string, expired bool) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	exp := time.Now().Add(time.Hour).Unix()
	if expired {
		exp = time.Now().Add(-time.Hour).Unix()
	}
	payloadBytes, _ := json.Marshal(map[string]interface{}{"sub": "admin", "exp": exp})
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(header + "." + payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return header + "." + payload + "." + sig
}

func TestHealthz(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetAggregatedHomepage(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/homepage", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp models.HomepageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Hero == nil {
		t.Error("expected hero in aggregated response")
	}
	if len(resp.Announcements) == 0 {
		t.Error("expected announcements in aggregated response")
	}
}

func TestListAnnouncements_WithLimitQueryParam(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/homepage/announcements?limit=1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var items []models.Announcement
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 announcement, got %d", len(items))
	}
}

func TestAdminCreateHero_RequiresAuth(t *testing.T) {
	mux := newTestMux()
	body, _ := json.Marshal(models.HeroBannerInput{Headline: "H", ImageURL: "https://cdn.example/x.jpg"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/homepage/hero", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}
}

func TestAdminCreateHero_RejectsExpiredToken(t *testing.T) {
	mux := newTestMux()
	body, _ := json.Marshal(models.HeroBannerInput{Headline: "H", ImageURL: "https://cdn.example/x.jpg"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/homepage/hero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+makeTestJWT(testJWTSecret, true))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", rec.Code)
	}
}

func TestAdminCreateHero_SucceedsWithValidToken(t *testing.T) {
	mux := newTestMux()
	body, _ := json.Marshal(models.HeroBannerInput{
		Headline: "Grand Reopening",
		ImageURL: "https://cdn.example/reopen.jpg",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/homepage/hero", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+makeTestJWT(testJWTSecret, false))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", rec.Code, rec.Body.String())
	}

	// Confirm the public hero endpoint now reflects the update.
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/homepage/hero", nil)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	var hero models.HeroBanner
	if err := json.Unmarshal(getRec.Body.Bytes(), &hero); err != nil {
		t.Fatalf("failed to decode hero: %v", err)
	}
	if hero.Headline != "Grand Reopening" {
		t.Errorf("expected updated headline, got %q", hero.Headline)
	}
}

func TestAdminDeleteAnnouncement_NotFoundReturns404(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/homepage/announcements/does-not-exist", nil)
	req.Header.Set("Authorization", "Bearer "+makeTestJWT(testJWTSecret, false))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListUpcomingEvents_OnlyFutureEvents(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/homepage/events", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var events []models.Event
	if err := json.Unmarshal(rec.Body.Bytes(), &events); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	now := time.Now().UTC()
	for _, e := range events {
		if e.StartDate.Before(now) {
			t.Errorf("expected only future events, got past event %q", e.Name)
		}
	}
}
