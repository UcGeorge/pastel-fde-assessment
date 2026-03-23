package config

import "os"

// Config holds all application configuration values.
type Config struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	AMLBaseURL string
	Port       string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		APIKey:     envOrDefault("SIGMA_API_KEY", "59d01b8c-7d29-4e25-8d93-c5c9f2bb0fdf"),
		APISecret:  envOrDefault("SIGMA_API_SECRET", "ad8a3c3f-1869-4947-94ad-bb216b792824"),
		BaseURL:    envOrDefault("SIGMA_BASE_URL", "https://sigmaprod.sabipay.com/"),
		AMLBaseURL: envOrDefault("SIGMA_AML_BASE_URL", "https://sigmaaml.sabipay.com/"),
		Port:       envOrDefault("PORT", "8080"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
