package repository

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/genoagolfclub/homepage-service/internal/models"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("resource not found")

// HomepageRepository is the storage abstraction for all homepage content.
//
// The production implementation for this deliverable is InMemoryHomepageRepository.
// A MongoDB-backed implementation follows the same interface and can be swapped
// in at wiring time (see cmd/server/main.go) without touching the service or
// handler layers, matching the existing Genoa Golf Club backend's Mongo-based
// repository pattern.
type HomepageRepository interface {
	GetActiveHero() (*models.HeroBanner, error)
	UpsertHero(banner *models.HeroBanner) (*models.HeroBanner, error)

	ListAnnouncements(limit int) ([]models.Announcement, error)
	CreateAnnouncement(a *models.Announcement) (*models.Announcement, error)
	DeleteAnnouncement(id string) error

	ListUpcomingEvents(limit int) ([]models.Event, error)

	ListMembershipHighlights() ([]models.MembershipHighlight, error)

	ListGalleryImages() ([]models.GalleryImage, error)
}

// InMemoryHomepageRepository is a thread-safe, in-memory implementation of
// HomepageRepository. It ships with representative seed data so the service
// is immediately usable for local development, demos, and automated tests.
type InMemoryHomepageRepository struct {
	mu            sync.RWMutex
	hero          *models.HeroBanner
	announcements map[string]models.Announcement
	events        map[string]models.Event
	membership    map[string]models.MembershipHighlight
	gallery       map[string]models.GalleryImage
	nextID        int
}

// NewInMemoryHomepageRepository builds a repository pre-populated with seed data.
func NewInMemoryHomepageRepository() *InMemoryHomepageRepository {
	repo := &InMemoryHomepageRepository{
		announcements: make(map[string]models.Announcement),
		events:        make(map[string]models.Event),
		membership:    make(map[string]models.MembershipHighlight),
		gallery:       make(map[string]models.GalleryImage),
	}
	repo.seed()
	return repo
}

func (r *InMemoryHomepageRepository) newID() string {
	r.nextID++
	return time.Now().UTC().Format("20060102150405") + "-" + itoa(r.nextID)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		return "-" + string(b)
	}
	return string(b)
}

func (r *InMemoryHomepageRepository) seed() {
	now := time.Now().UTC()

	r.hero = &models.HeroBanner{
		ID:        "hero-1",
		Headline:  "Welcome to Genoa Golf Club",
		Subtext:   "Championship greens. Timeless hospitality.",
		ImageURL:  "https://cdn.genoagolfclub.example/hero/clubhouse.jpg",
		CTALabel:  "Book a Tee Time",
		CTAUrl:    "/tee-times",
		Active:    true,
		UpdatedAt: now,
	}

	seedAnnouncements := []models.Announcement{
		{ID: "ann-1", Title: "Fall Membership Drive Open", Body: "New members save 15% through October.", PostedAt: now.Add(-24 * time.Hour)},
		{ID: "ann-2", Title: "Clubhouse Renovation Complete", Body: "The new grill room is now open daily.", PostedAt: now.Add(-72 * time.Hour)},
		{ID: "ann-3", Title: "Course Aeration Schedule", Body: "Greens aeration runs the first week of next month.", PostedAt: now.Add(-96 * time.Hour)},
	}
	for _, a := range seedAnnouncements {
		r.announcements[a.ID] = a
	}

	seedEvents := []models.Event{
		{ID: "evt-1", Name: "Member-Guest Invitational", StartDate: now.Add(7 * 24 * time.Hour), Location: "Championship Course"},
		{ID: "evt-2", Name: "Junior Golf Clinic", StartDate: now.Add(3 * 24 * time.Hour), Location: "Practice Range"},
		{ID: "evt-3", Name: "Autumn Charity Scramble", StartDate: now.Add(21 * 24 * time.Hour), Location: "Championship Course"},
		{ID: "evt-4", Name: "Ladies' Nine & Dine", StartDate: now.Add(-2 * 24 * time.Hour), Location: "North Nine"}, // past, should be filtered
	}
	for _, e := range seedEvents {
		r.events[e.ID] = e
	}

	seedMembership := []models.MembershipHighlight{
		{ID: "mem-1", Tier: "Individual", Price: 249, Benefits: []string{"Unlimited weekday play", "10% pro shop discount"}},
		{ID: "mem-2", Tier: "Family", Price: 399, Benefits: []string{"Unlimited play for household", "Free junior clinics"}},
		{ID: "mem-3", Tier: "Corporate", Price: 899, Benefits: []string{"4 transferable memberships", "Private event hosting"}},
	}
	for _, m := range seedMembership {
		r.membership[m.ID] = m
	}

	seedGallery := []models.GalleryImage{
		{ID: "img-1", URL: "https://cdn.genoagolfclub.example/gallery/18th-green.jpg", Caption: "The 18th green at sunset"},
		{ID: "img-2", URL: "https://cdn.genoagolfclub.example/gallery/clubhouse.jpg", Caption: "The renovated clubhouse"},
		{ID: "img-3", URL: "https://cdn.genoagolfclub.example/gallery/tournament.jpg", Caption: "2025 Member-Guest Invitational"},
	}
	for _, g := range seedGallery {
		r.gallery[g.ID] = g
	}
}

func (r *InMemoryHomepageRepository) GetActiveHero() (*models.HeroBanner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.hero == nil {
		return nil, ErrNotFound
	}
	cp := *r.hero
	return &cp, nil
}

func (r *InMemoryHomepageRepository) UpsertHero(banner *models.HeroBanner) (*models.HeroBanner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if banner.ID == "" {
		banner.ID = r.newID()
	}
	banner.UpdatedAt = time.Now().UTC()
	banner.Active = true
	cp := *banner
	r.hero = &cp
	out := *r.hero
	return &out, nil
}

func (r *InMemoryHomepageRepository) ListAnnouncements(limit int) ([]models.Announcement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Announcement, 0, len(r.announcements))
	for _, a := range r.announcements {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PostedAt.After(out[j].PostedAt) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *InMemoryHomepageRepository) CreateAnnouncement(a *models.Announcement) (*models.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a.ID == "" {
		a.ID = r.newID()
	}
	if a.PostedAt.IsZero() {
		a.PostedAt = time.Now().UTC()
	}
	r.announcements[a.ID] = *a
	out := *a
	return &out, nil
}

func (r *InMemoryHomepageRepository) DeleteAnnouncement(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.announcements[id]; !ok {
		return ErrNotFound
	}
	delete(r.announcements, id)
	return nil
}

func (r *InMemoryHomepageRepository) ListUpcomingEvents(limit int) ([]models.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now().UTC()
	out := make([]models.Event, 0, len(r.events))
	for _, e := range r.events {
		if e.StartDate.After(now) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartDate.Before(out[j].StartDate) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *InMemoryHomepageRepository) ListMembershipHighlights() ([]models.MembershipHighlight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.MembershipHighlight, 0, len(r.membership))
	for _, m := range r.membership {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Price < out[j].Price })
	return out, nil
}

func (r *InMemoryHomepageRepository) ListGalleryImages() ([]models.GalleryImage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.GalleryImage, 0, len(r.gallery))
	for _, g := range r.gallery {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
