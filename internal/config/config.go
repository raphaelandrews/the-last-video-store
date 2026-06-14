package config

import (
	"os"
	"time"
)

type Config struct {
	DBPath      string
	JWTSecret   string
	AESKey      string
	ServerPort  string
	APIBaseURL  string
	Environment string
}

func Load() *Config {
	return &Config{
		DBPath:      envOrDefault("TLVS_DB_PATH", "thelastvideostore.db"),
		JWTSecret:   envOrDefault("TLVS_JWT_SECRET", "change-me-in-production"),
		AESKey:      envOrDefault("TLVS_AES_KEY", "0123456789abcdef0123456789abcdef"),
		ServerPort:  envOrDefault("TLVS_SERVER_PORT", "8080"),
		APIBaseURL:  envOrDefault("TLVS_API_BASE_URL", "http://localhost:8080"),
		Environment: envOrDefault("TLVS_ENV", "development"),
	}
}

func MustLoad() *Config {
	required := map[string]string{
		"TLVS_JWT_SECRET": "",
		"TLVS_AES_KEY":    "",
	}
	for key := range required {
		if os.Getenv(key) == "" {
			panic("missing required environment variable: " + key)
		}
	}
	return Load()
}

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
	LockoutDuration      = 30 * time.Minute
	MaxLoginAttempts     = 5
	BcryptCost           = 12
	RateLimitPerMinute   = 100
	DefaultPageSize      = 20
	MaxSearchResults     = 10

	FormatVHS    = "VHS"
	FormatDVD    = "DVD"
	FormatBluRay = "Blu-ray"

	VHSLateFeeRate = 2.00
	DVDLateFeeRate = 3.00
	RewindFeeCost  = 1.00
)

var GenreList = []string{"Action", "Comedy", "Horror", "SciFi", "Drama", "Thriller", "Romance", "Animation"}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
