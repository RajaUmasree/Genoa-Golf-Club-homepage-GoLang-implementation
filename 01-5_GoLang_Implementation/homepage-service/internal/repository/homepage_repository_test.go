package repository

import (
	"testing"
	"time"

	"github.com/genoagolfclub/homepage-service/internal/models"
)

func TestGetActiveHero_ReturnsSeededHero(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	hero, err := repo.GetActiveHero()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hero.Headline == "" {
		t.Error("expected seeded hero to have a headline")
	}
}

func TestUpsertHero_UpdatesActiveHero(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	updated, err := repo.UpsertHero(&models.HeroBanner{
		Headline: "New Season, New Course",
		ImageURL: "https://cdn.example/new.jpg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID == "" {
		t.Error("expected a generated ID")
	}

	got, err := repo.GetActiveHero()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Headline != "New Season, New Course" {
		t.Errorf("expected updated headline, got %q", got.Headline)
	}
}

func TestListAnnouncements_RespectsLimitAndOrder(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	items, err := repo.ListAnnouncements(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if !items[0].PostedAt.After(items[1].PostedAt) {
		t.Error("expected announcements sorted newest first")
	}
}

func TestCreateAndDeleteAnnouncement(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	created, err := repo.CreateAnnouncement(&models.Announcement{Title: "Test", Body: "Body"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated ID")
	}

	if err := repo.DeleteAnnouncement(created.ID); err != nil {
		t.Fatalf("unexpected error deleting: %v", err)
	}

	if err := repo.DeleteAnnouncement(created.ID); err != ErrNotFound {
		t.Errorf("expected ErrNotFound on second delete, got %v", err)
	}
}

func TestListUpcomingEvents_ExcludesPastEvents(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	events, err := repo.ListUpcomingEvents(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Now().UTC()
	for _, e := range events {
		if e.StartDate.Before(now) {
			t.Errorf("expected only future events, found past event %q", e.Name)
		}
	}
	// Seed data includes one past event ("Ladies' Nine & Dine") that must be filtered out.
	for _, e := range events {
		if e.Name == "Ladies' Nine & Dine" {
			t.Error("past event should have been filtered out")
		}
	}
}

func TestListUpcomingEvents_SortedAscendingAndLimited(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	events, err := repo.ListUpcomingEvents(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].StartDate.After(events[1].StartDate) {
		t.Error("expected events sorted soonest-first")
	}
}

func TestListMembershipHighlights_SortedByPrice(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	items, err := repo.ListMembershipHighlights()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(items); i++ {
		if items[i].Price < items[i-1].Price {
			t.Error("expected membership highlights sorted ascending by price")
		}
	}
}

func TestListGalleryImages_ReturnsSeeded(t *testing.T) {
	repo := NewInMemoryHomepageRepository()
	items, err := repo.ListGalleryImages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) == 0 {
		t.Error("expected seeded gallery images")
	}
}
