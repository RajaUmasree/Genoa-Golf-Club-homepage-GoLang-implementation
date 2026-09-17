package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/genoagolfclub/homepage-service/internal/models"
	"github.com/genoagolfclub/homepage-service/internal/repository"
	"github.com/genoagolfclub/homepage-service/internal/service"
)

// HomepageHandler holds the http.HandlerFunc implementations for every
// homepage endpoint. It depends only on the service layer, so it is easy to
// unit test with httptest and an in-memory repository.
type HomepageHandler struct {
	svc *service.HomepageService
}

func NewHomepageHandler(svc *service.HomepageService) *HomepageHandler {
	return &HomepageHandler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func (h *HomepageHandler) respondErr(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, models.NewError(models.CodeNotFound, "resource not found"))
		return
	}
	writeJSON(w, http.StatusInternalServerError, models.NewError(models.CodeInternal, "internal server error"))
}

func limitFromQuery(r *http.Request, def int) int {
	v := r.URL.Query().Get("limit")
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

// PathValue helper — wraps r.PathValue (Go 1.22 ServeMux) so handlers stay
// testable without booting a full mux for unit tests that set it manually.
func pathID(r *http.Request) string { return r.PathValue("id") }

// GetAggregatedHomepage handles GET /api/v1/homepage
func (h *HomepageHandler) GetAggregatedHomepage(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetAggregatedHomepage()
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetHero handles GET /api/v1/homepage/hero
func (h *HomepageHandler) GetHero(w http.ResponseWriter, r *http.Request) {
	hero, err := h.svc.GetHero()
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, hero)
}

// CreateHero handles POST /api/v1/admin/homepage/hero
func (h *HomepageHandler) CreateHero(w http.ResponseWriter, r *http.Request) {
	var input models.HeroBannerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, models.NewError(models.CodeValidation, "invalid request body"))
		return
	}
	if input.Headline == "" || input.ImageURL == "" {
		writeJSON(w, http.StatusBadRequest, models.NewError(models.CodeValidation, "headline and imageUrl are required"))
		return
	}
	hero, err := h.svc.SetHero(input)
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, hero)
}

// UpdateHero handles PUT /api/v1/admin/homepage/hero/{id}
func (h *HomepageHandler) UpdateHero(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	var input models.HeroBannerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, models.NewError(models.CodeValidation, "invalid request body"))
		return
	}
	hero, err := h.svc.UpdateHero(id, input)
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, hero)
}

// ListAnnouncements handles GET /api/v1/homepage/announcements
func (h *HomepageHandler) ListAnnouncements(w http.ResponseWriter, r *http.Request) {
	limit := limitFromQuery(r, 3)
	items, err := h.svc.ListAnnouncements(limit)
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// CreateAnnouncement handles POST /api/v1/admin/homepage/announcements
func (h *HomepageHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var input models.AnnouncementInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, models.NewError(models.CodeValidation, "invalid request body"))
		return
	}
	if input.Title == "" || input.Body == "" {
		writeJSON(w, http.StatusBadRequest, models.NewError(models.CodeValidation, "title and body are required"))
		return
	}
	a, err := h.svc.CreateAnnouncement(input)
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

// DeleteAnnouncement handles DELETE /api/v1/admin/homepage/announcements/{id}
func (h *HomepageHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	if err := h.svc.DeleteAnnouncement(id); err != nil {
		h.respondErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListUpcomingEvents handles GET /api/v1/homepage/events
func (h *HomepageHandler) ListUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	limit := limitFromQuery(r, 5)
	items, err := h.svc.ListUpcomingEvents(limit)
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// ListMembershipHighlights handles GET /api/v1/homepage/membership-highlights
func (h *HomepageHandler) ListMembershipHighlights(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListMembershipHighlights()
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// ListGallery handles GET /api/v1/homepage/gallery
func (h *HomepageHandler) ListGallery(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListGallery()
	if err != nil {
		h.respondErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// Healthz handles GET /healthz
func (h *HomepageHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
