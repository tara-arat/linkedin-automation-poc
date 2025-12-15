package models

import "time"

// Config represents the application configuration
type Config struct {
	LinkedIn    LinkedInConfig    `yaml:"linkedin" json:"linkedin"`
	Stealth     StealthConfig     `yaml:"stealth" json:"stealth"`
	RateLimits  RateLimitsConfig  `yaml:"rate_limits" json:"rate_limits"`
	Storage     StorageConfig     `yaml:"storage" json:"storage"`
	Logging     LoggingConfig     `yaml:"logging" json:"logging"`
}

// LinkedInConfig contains LinkedIn authentication settings
type LinkedInConfig struct {
	Email           string `yaml:"email" json:"email"`
	Password        string `yaml:"password" json:"password"`
	SessionPath     string `yaml:"session_path" json:"session_path"`
	BaseURL         string `yaml:"base_url" json:"base_url"`
}

// StealthConfig contains anti-detection settings
type StealthConfig struct {
	EnableMouseMovement    bool          `yaml:"enable_mouse_movement" json:"enable_mouse_movement"`
	EnableTypingSimulation bool          `yaml:"enable_typing_simulation" json:"enable_typing_simulation"`
	EnableRandomScrolling  bool          `yaml:"enable_random_scrolling" json:"enable_random_scrolling"`
	EnableHovering         bool          `yaml:"enable_hovering" json:"enable_hovering"`
	MinActionDelay         time.Duration `yaml:"min_action_delay" json:"min_action_delay"`
	MaxActionDelay         time.Duration `yaml:"max_action_delay" json:"max_action_delay"`
	BusinessHoursOnly      bool          `yaml:"business_hours_only" json:"business_hours_only"`
	BusinessHoursStart     int           `yaml:"business_hours_start" json:"business_hours_start"`
	BusinessHoursEnd       int           `yaml:"business_hours_end" json:"business_hours_end"`
}

// RateLimitsConfig defines rate limiting settings
type RateLimitsConfig struct {
	MaxConnectionsPerDay int           `yaml:"max_connections_per_day" json:"max_connections_per_day"`
	MaxMessagesPerDay    int           `yaml:"max_messages_per_day" json:"max_messages_per_day"`
	MaxSearchesPerHour   int           `yaml:"max_searches_per_hour" json:"max_searches_per_hour"`
	CooldownPeriod       time.Duration `yaml:"cooldown_period" json:"cooldown_period"`
}

// StorageConfig defines database settings
type StorageConfig struct {
	DatabasePath string `yaml:"database_path" json:"database_path"`
	BackupPath   string `yaml:"backup_path" json:"backup_path"`
}

// LoggingConfig defines logging settings
type LoggingConfig struct {
	Level      string `yaml:"level" json:"level"`
	Format     string `yaml:"format" json:"format"`
	OutputPath string `yaml:"output_path" json:"output_path"`
}
