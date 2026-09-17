package service

import (
	"github.com/genoagolfclub/homepage-service/internal/models"
	"github.com/genoagolfclub/homepage-service/internal/repository"
)

const (
	defaultAnnouncementLimit = 3
	defaultEventLimit        = 5
)

// HomepageService contains the business rules for building homepage content.
// It is intentionally storage-agnostic — it depends only on the
// repository.HomepageRepository interface.
type HomepageService struct {
	repo repository.HomepageRepository
}

func NewHomepageService(repo repository.HomepageRepository) *HomepageService {
	return &HomepageService{repo: repo}
}

func (s *HomepageService) GetHero() (*models.HeroBanner, error) {
	return s.repo.GetActiveHero()
}

func (s *HomepageService) SetHero(input models.HeroBannerInput) (*models.HeroBanner, error) {
	banner := &models.HeroBanner{
		Headline: input.Headline,
		Subtext:  input.Subtext,
		ImageURL: input.ImageURL,
		CTALabel: input.CTALabel,
		CTAUrl:   input.CTAUrl,
	}
	return s.repo.UpsertHero(banner)
}

func (s *HomepageService) UpdateHero(id string, input models.HeroBannerInput) (*models.HeroBanner, error) {
	existing, err := s.repo.GetActiveHero()
	if err != nil {
		return nil, err
	}
	if existing.ID != id {
		return nil, repository.ErrNotFound
	}
	existing.Headline = input.Headline
	existing.Subtext = input.Subtext
	existing.ImageURL = input.ImageURL
	existing.CTALabel = input.CTALabel
	existing.CTAUrl = input.CTAUrl
	return s.repo.UpsertHero(existing)
}

func (s *HomepageService) ListAnnouncements(limit int) ([]models.Announcement, error) {
	if limit <= 0 {
		limit = defaultAnnouncementLimit
	}
	return s.repo.ListAnnouncements(limit)
}

func (s *HomepageService) CreateAnnouncement(input models.AnnouncementInput) (*models.Announcement, error) {
	a := &models.Announcement{Title: input.Title, Body: input.Body}
	return s.repo.CreateAnnouncement(a)
}

func (s *HomepageService) DeleteAnnouncement(id string) error {
	return s.repo.DeleteAnnouncement(id)
}

func (s *HomepageService) ListUpcomingEvents(limit int) ([]models.Event, error) {
	if limit <= 0 {
		limit = defaultEventLimit
	}
	return s.repo.ListUpcomingEvents(limit)
}

func (s *HomepageService) ListMembershipHighlights() ([]models.MembershipHighlight, error) {
	return s.repo.ListMembershipHighlights()
}

func (s *HomepageService) ListGallery() ([]models.GalleryImage, error) {
	return s.repo.ListGalleryImages()
}

// GetAggregatedHomepage fans out to every section and combines the results
// into a single HomepageResponse for the landing page's single-call use case.
func (s *HomepageService) GetAggregatedHomepage() (*models.HomepageResponse, error) {
	resp := &models.HomepageResponse{}

	if hero, err := s.repo.GetActiveHero(); err == nil {
		resp.Hero = hero
	} else if err != repository.ErrNotFound {
		return nil, err
	}

	announcements, err := s.repo.ListAnnouncements(defaultAnnouncementLimit)
	if err != nil {
		return nil, err
	}
	resp.Announcements = announcements

	events, err := s.repo.ListUpcomingEvents(defaultEventLimit)
	if err != nil {
		return nil, err
	}
	resp.UpcomingEvents = events

	membership, err := s.repo.ListMembershipHighlights()
	if err != nil {
		return nil, err
	}
	resp.MembershipHighlights = membership

	gallery, err := s.repo.ListGalleryImages()
	if err != nil {
		return nil, err
	}
	resp.Gallery = gallery

	return resp, nil
}
