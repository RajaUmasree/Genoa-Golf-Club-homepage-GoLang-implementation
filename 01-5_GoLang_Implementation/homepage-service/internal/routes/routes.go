package routes

import (
	"net/http"

	"github.com/genoagolfclub/homepage-service/internal/handlers"
	"github.com/genoagolfclub/homepage-service/internal/middleware"
)

// RegisterRoutes wires the public and admin homepage routes onto the given
// http.ServeMux using Go 1.22's method+path pattern matching
// (e.g. "GET /api/v1/homepage/hero").
func RegisterRoutes(mux *http.ServeMux, h *handlers.HomepageHandler, jwtSecret string) {
	mux.HandleFunc("GET /healthz", h.Healthz)

	// Public, read-only homepage routes.
	mux.HandleFunc("GET /api/v1/homepage", h.GetAggregatedHomepage)
	mux.HandleFunc("GET /api/v1/homepage/hero", h.GetHero)
	mux.HandleFunc("GET /api/v1/homepage/announcements", h.ListAnnouncements)
	mux.HandleFunc("GET /api/v1/homepage/events", h.ListUpcomingEvents)
	mux.HandleFunc("GET /api/v1/homepage/membership-highlights", h.ListMembershipHighlights)
	mux.HandleFunc("GET /api/v1/homepage/gallery", h.ListGallery)

	// Admin, JWT-protected write routes.
	mux.Handle("POST /api/v1/admin/homepage/hero",
		middleware.JWTAuth(jwtSecret, http.HandlerFunc(h.CreateHero)))
	mux.Handle("PUT /api/v1/admin/homepage/hero/{id}",
		middleware.JWTAuth(jwtSecret, http.HandlerFunc(h.UpdateHero)))
	mux.Handle("POST /api/v1/admin/homepage/announcements",
		middleware.JWTAuth(jwtSecret, http.HandlerFunc(h.CreateAnnouncement)))
	mux.Handle("DELETE /api/v1/admin/homepage/announcements/{id}",
		middleware.JWTAuth(jwtSecret, http.HandlerFunc(h.DeleteAnnouncement)))
}
