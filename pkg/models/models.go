package models

import "time"

// Profile represents a LinkedIn profile
type Profile struct {
	ID          int64     `json:"id"`
	ProfileURL  string    `json:"profile_url"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	DiscoveredAt time.Time `json:"discovered_at"`
}

// ConnectionRequest represents a sent connection request
type ConnectionRequest struct {
	ID          int64     `json:"id"`
	ProfileURL  string    `json:"profile_url"`
	ProfileName string    `json:"profile_name"`
	Message     string    `json:"message"`
	SentAt      time.Time `json:"sent_at"`
	Status      string    `json:"status"` // pending, accepted, rejected
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
}

// Message represents a sent message
type Message struct {
	ID          int64     `json:"id"`
	ProfileURL  string    `json:"profile_url"`
	ProfileName string    `json:"profile_name"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
	IsFollowUp  bool      `json:"is_follow_up"`
}

// SearchCriteria represents search parameters
type SearchCriteria struct {
	Keywords  []string `json:"keywords"`
	JobTitle  string   `json:"job_title"`
	Company   string   `json:"company"`
	Location  string   `json:"location"`
	MaxResults int     `json:"max_results"`
}

// ActivityStats tracks daily activity
type ActivityStats struct {
	Date              string `json:"date"`
	ConnectionsSent   int    `json:"connections_sent"`
	MessagesSent      int    `json:"messages_sent"`
	SearchesPerformed int    `json:"searches_performed"`
	ProfilesViewed    int    `json:"profiles_viewed"`
}
