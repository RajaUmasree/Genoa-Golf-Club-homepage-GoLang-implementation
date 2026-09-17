//go:build mongo

// Package repository — MongoDB-backed implementation.
//
// This file is excluded from the default build (build tag "mongo") so the
// deliverable compiles and runs with zero external dependencies out of the
// box. It documents the intended production wiring against the existing
// Genoa Golf Club MongoDB cluster and matching driver version used elsewhere
// in the GGC backend.
//
// To build with this implementation instead of the in-memory one:
//
//	go build -tags mongo ./...
//
// and set MONGO_URI / MONGO_DB in the environment (see internal/config).
package repository

import (
	"context"
	"time"

	"github.com/genoagolfclub/homepage-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NOTE: enabling this file (via `-tags mongo`) requires adding
// go.mongodb.org/mongo-driver to go.mod, e.g.:
//   go get go.mongodb.org/mongo-driver@v1.16.0
// It is intentionally excluded from the default module requirements so the
// service builds and tests offline with zero external dependencies.

type MongoHomepageRepository struct {
	db *mongo.Database
}

func NewMongoHomepageRepository(db *mongo.Database) *MongoHomepageRepository {
	return &MongoHomepageRepository{db: db}
}

func (r *MongoHomepageRepository) heroCol() *mongo.Collection          { return r.db.Collection("homepage_hero") }
func (r *MongoHomepageRepository) announceCol() *mongo.Collection     { return r.db.Collection("homepage_announcements") }
func (r *MongoHomepageRepository) eventsCol() *mongo.Collection       { return r.db.Collection("events") }
func (r *MongoHomepageRepository) membershipCol() *mongo.Collection   { return r.db.Collection("membership_tiers") }
func (r *MongoHomepageRepository) galleryCol() *mongo.Collection      { return r.db.Collection("homepage_gallery") }

func (r *MongoHomepageRepository) GetActiveHero() (*models.HeroBanner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var hero models.HeroBanner
	err := r.heroCol().FindOne(ctx, bson.M{"active": true}).Decode(&hero)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &hero, nil
}

func (r *MongoHomepageRepository) UpsertHero(banner *models.HeroBanner) (*models.HeroBanner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	banner.UpdatedAt = time.Now().UTC()
	banner.Active = true
	_, err := r.heroCol().UpdateOne(ctx, bson.M{"_id": banner.ID},
		bson.M{"$set": banner}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return banner, nil
}

// ListAnnouncements, CreateAnnouncement, DeleteAnnouncement, ListUpcomingEvents,
// ListMembershipHighlights, and ListGalleryImages follow the same pattern:
// standard Mongo find/sort/limit or insert/delete calls against the
// collections above. Omitted here for brevity — implement following the
// existing GGC backend's Mongo repository conventions before enabling the
// "mongo" build tag in production (see 01-8 release readiness runbook).
