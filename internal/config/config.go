package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/keerthana/linkedin-automation-poc/pkg/models"
	"gopkg.in/yaml.v3"
)

// Load reads configuration from file and environment variables
func Load(configPath string) (*models.Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Read config file
	cfg, err := loadFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Override with environment variables
	overrideFromEnv(cfg)

	// Validate configuration
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// loadFromFile reads configuration from YAML file
func loadFromFile(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file doesn't exist
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// overrideFromEnv overrides config values with environment variables
func overrideFromEnv(cfg *models.Config) {
	if email := os.Getenv("LINKEDIN_EMAIL"); email != "" {
		cfg.LinkedIn.Email = email
	}
	if password := os.Getenv("LINKEDIN_PASSWORD"); password != "" {
		cfg.LinkedIn.Password = password
	}
	if sessionPath := os.Getenv("LINKEDIN_SESSION_PATH"); sessionPath != "" {
		cfg.LinkedIn.SessionPath = sessionPath
	}
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.Storage.DatabasePath = dbPath
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Logging.Level = logLevel
	}
}

// validate ensures all required configuration values are present
func validate(cfg *models.Config) error {
	if cfg.LinkedIn.Email == "" {
		return fmt.Errorf("linkedin email is required")
	}
	if cfg.LinkedIn.Password == "" {
		return fmt.Errorf("linkedin password is required")
	}
	if cfg.RateLimits.MaxConnectionsPerDay <= 0 {
		return fmt.Errorf("max connections per day must be positive")
	}
	if cfg.Stealth.MinActionDelay > cfg.Stealth.MaxActionDelay {
		return fmt.Errorf("min action delay cannot be greater than max action delay")
	}
	return nil
}

// defaultConfig returns a configuration with sensible defaults
func defaultConfig() *models.Config {
	return &models.Config{
		LinkedIn: models.LinkedInConfig{
			SessionPath: "./sessions",
			BaseURL:     "https://www.linkedin.com",
		},
		Stealth: models.StealthConfig{
			EnableMouseMovement:    true,
			EnableTypingSimulation: true,
			EnableRandomScrolling:  true,
			EnableHovering:         true,
			MinActionDelay:         2 * time.Second,
			MaxActionDelay:         5 * time.Second,
			BusinessHoursOnly:      true,
			BusinessHoursStart:     9,  // 9 AM
			BusinessHoursEnd:       17, // 5 PM
		},
		RateLimits: models.RateLimitsConfig{
			MaxConnectionsPerDay: 20,
			MaxMessagesPerDay:    15,
			MaxSearchesPerHour:   5,
			CooldownPeriod:       30 * time.Minute,
		},
		Storage: models.StorageConfig{
			DatabasePath: "./data/automation.db",
			BackupPath:   "./data/backups",
		},
		Logging: models.LoggingConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "./logs/automation.log",
		},
	}
}
