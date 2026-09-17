package config

import "os"

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	Port      string
	JWTSecret string
	MongoURI  string // used only when built with -tags mongo
	MongoDB   string // used only when built with -tags mongo
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Load reads configuration from the environment, applying sane local
// development defaults.
func Load() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "dev-only-secret-change-me"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DB", "genoa_golf_club"),
	}
}
