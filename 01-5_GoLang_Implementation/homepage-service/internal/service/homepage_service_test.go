package service

import (
	"testing"

	"github.com/genoagolfclub/homepage-service/internal/models"
	"github.com/genoagolfclub/homepage-service/internal/repository"
)

func newTestService() *HomepageService {
	repo := repository.NewInMemoryHomepageRepository()
	return NewHomepageService(repo)
}

func TestGetAggregatedHomepage_PopulatesAllSections(t *testing.T) {
	svc := newTestService()
	resp, err := svc.GetAggregatedHomepage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Hero == nil {
		t.Error("expected hero to be populated")
	}
	if len(resp.Announcements) == 0 {
		t.Error("expected announcements to be populated")
	}
	if len(resp.UpcomingEvents) == 0 {
		t.Error("expected upcoming events to be populated")
	}
	if len(resp.MembershipHighlights) == 0 {
		t.Error("expected membership highlights to be populated")
	}
	if len(resp.Gallery) == 0 {
		t.Error("expected gallery to be populated")
	}
}

func TestListAnnouncements_DefaultsLimitWhenZero(t *testing.T) {
	svc := newTestService()
	items, err := svc.ListAnnouncements(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != defaultAnnouncementLimit {
		t.Errorf("expected default limit %d, got %d", defaultAnnouncementLimit, len(items))
	}
}

func TestListUpcomingEvents_DefaultsLimitWhenZero(t *testing.T) {
	svc := newTestService()
	items, err := svc.ListUpcomingEvents(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) > defaultEventLimit {
		t.Errorf("expected at most %d events, got %d", defaultEventLimit, len(items))
	}
}

func TestSetHero_CreatesActiveHero(t *testing.T) {
	svc := newTestService()
	hero, err := svc.SetHero(models.HeroBannerInput{
		Headline: "Grand Reopening",
		ImageURL: "https://cdn.example/reopen.jpg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hero.Active {
		t.Error("expected new hero to be active")
	}

	got, err := svc.GetHero()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Headline != "Grand Reopening" {
		t.Errorf("expected updated headline, got %q", got.Headline)
	}
}

func TestUpdateHero_NotFoundWhenIDMismatch(t *testing.T) {
	svc := newTestService()
	_, err := svc.UpdateHero("does-not-exist", models.HeroBannerInput{
		Headline: "X",
		ImageURL: "https://cdn.example/x.jpg",
	})
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateAndDeleteAnnouncement(t *testing.T) {
	svc := newTestService()
	created, err := svc.CreateAnnouncement(models.AnnouncementInput{Title: "T", Body: "B"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.DeleteAnnouncement(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.DeleteAnnouncement(created.ID); err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
