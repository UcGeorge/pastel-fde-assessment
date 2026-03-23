package config

import "os"

// Config holds all application configuration values.
type Config struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	AMLBaseURL string
	Port       string
	UseMock    bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		APIKey:     envOrDefault("SIGMA_API_KEY", "59d01b8c-..."),
		APISecret:  envOrDefault("SIGMA_API_SECRET", "ad8a3c3f-..."),
		BaseURL:    envOrDefault("SIGMA_BASE_URL", "https://sigmaprod.sabipay.com/"),
		AMLBaseURL: envOrDefault("SIGMA_AML_BASE_URL", "https://sigmaaml.sabipay.com/"),
		Port:       envOrDefault("PORT", "80"),
		UseMock:    envOrDefault("USE_MOCK", "true") == "true",
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
