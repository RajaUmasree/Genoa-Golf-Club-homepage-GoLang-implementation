package models

import "time"

// HeroBanner represents the main marketing banner shown at the top of the homepage.
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

// Announcement is a short club news item shown on the homepage.
type Announcement struct {
	ID       string    `json:"id" bson:"_id"`
	Title    string    `json:"title" bson:"title"`
	Body     string    `json:"body" bson:"body"`
	PostedAt time.Time `json:"postedAt" bson:"postedAt"`
}

// Event is an upcoming tournament / club event.
type Event struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	Location  string    `json:"location" bson:"location"`
}

// MembershipHighlight summarizes a membership tier for the homepage.
type MembershipHighlight struct {
	ID       string   `json:"id" bson:"_id"`
	Tier     string   `json:"tier" bson:"tier"`
	Price    float64  `json:"price" bson:"price"`
	Benefits []string `json:"benefits" bson:"benefits"`
}

// GalleryImage is a single curated image shown in the homepage gallery.
type GalleryImage struct {
	ID      string `json:"id" bson:"_id"`
	URL     string `json:"url" bson:"url"`
	Caption string `json:"caption" bson:"caption"`
}

// HomepageResponse aggregates every homepage section into a single payload.
type HomepageResponse struct {
	Hero                 *HeroBanner           `json:"hero"`
	Announcements        []Announcement        `json:"announcements"`
	UpcomingEvents       []Event               `json:"upcomingEvents"`
	MembershipHighlights []MembershipHighlight `json:"membershipHighlights"`
	Gallery              []GalleryImage        `json:"gallery"`
}

// HeroBannerInput is the writable subset of HeroBanner accepted from admin requests.
type HeroBannerInput struct {
	Headline string `json:"headline" binding:"required"`
	Subtext  string `json:"subtext"`
	ImageURL string `json:"imageUrl" binding:"required"`
	CTALabel string `json:"ctaLabel"`
	CTAUrl   string `json:"ctaUrl"`
}

// AnnouncementInput is the writable subset of Announcement accepted from admin requests.
type AnnouncementInput struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// ErrorResponse is the standard error envelope returned by every endpoint.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error codes used across the service.
const (
	CodeValidation   = "VALIDATION_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeNotFound     = "NOT_FOUND"
	CodeInternal     = "INTERNAL_ERROR"
)

func NewError(code, message string) ErrorResponse {
	return ErrorResponse{Error: ErrorBody{Code: code, Message: message}}
}
