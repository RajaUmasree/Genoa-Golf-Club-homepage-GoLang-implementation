// Command server runs the Genoa Golf Club Homepage service (GGC_GL_01).
package main

import (
	"log"
	"net/http"

	"github.com/genoagolfclub/homepage-service/internal/config"
	"github.com/genoagolfclub/homepage-service/internal/handlers"
	"github.com/genoagolfclub/homepage-service/internal/repository"
	"github.com/genoagolfclub/homepage-service/internal/routes"
	"github.com/genoagolfclub/homepage-service/internal/service"
)

func main() {
	cfg := config.Load()

	// Repository: in-memory by default. Build with `-tags mongo` and set
	// MONGO_URI/MONGO_DB to use the MongoDB-backed implementation instead
	// (see internal/repository/homepage_repository_mongo.go).
	repo := repository.NewInMemoryHomepageRepository()

	svc := service.NewHomepageService(repo)
	handler := handlers.NewHomepageHandler(svc)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, handler, cfg.JWTSecret)

	log.Printf("Genoa Golf Club Homepage service listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
