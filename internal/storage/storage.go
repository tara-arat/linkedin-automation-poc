package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/keerthana/linkedin-automation-poc/pkg/models"
)

// Storage handles database operations
type Storage struct {
	db *sql.DB
}

// New creates a new storage instance
func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{db: db}
	
	if err := storage.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return storage, nil
}

// initSchema creates database tables
func (s *Storage) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT UNIQUE NOT NULL,
		name TEXT,
		title TEXT,
		company TEXT,
		location TEXT,
		discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS connection_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT NOT NULL,
		profile_name TEXT,
		message TEXT,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		status TEXT DEFAULT 'pending',
		accepted_at DATETIME,
		UNIQUE(profile_url)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT NOT NULL,
		profile_name TEXT,
		content TEXT NOT NULL,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_follow_up BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS activity_stats (
		date TEXT PRIMARY KEY,
		connections_sent INTEGER DEFAULT 0,
		messages_sent INTEGER DEFAULT 0,
		searches_performed INTEGER DEFAULT 0,
		profiles_viewed INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_profiles_url ON profiles(profile_url);
	CREATE INDEX IF NOT EXISTS idx_connections_url ON connection_requests(profile_url);
	CREATE INDEX IF NOT EXISTS idx_connections_status ON connection_requests(status);
	CREATE INDEX IF NOT EXISTS idx_messages_url ON messages(profile_url);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveProfile saves a profile to the database
func (s *Storage) SaveProfile(profile *models.Profile) error {
	query := `
		INSERT OR IGNORE INTO profiles (profile_url, name, title, company, location, discovered_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, profile.ProfileURL, profile.Name, profile.Title,
		profile.Company, profile.Location, profile.DiscoveredAt)
	return err
}

// ProfileExists checks if a profile already exists
func (s *Storage) ProfileExists(profileURL string) (bool, error) {
	query := "SELECT COUNT(*) FROM profiles WHERE profile_url = ?"
	var count int
	err := s.db.QueryRow(query, profileURL).Scan(&count)
	return count > 0, err
}

// SaveConnectionRequest saves a connection request
func (s *Storage) SaveConnectionRequest(req *models.ConnectionRequest) error {
	query := `
		INSERT OR REPLACE INTO connection_requests (profile_url, profile_name, message, sent_at, status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, req.ProfileURL, req.ProfileName, req.Message, req.SentAt, req.Status)
	return err
}

// ConnectionRequestExists checks if a connection request was already sent
func (s *Storage) ConnectionRequestExists(profileURL string) (bool, error) {
	query := "SELECT COUNT(*) FROM connection_requests WHERE profile_url = ?"
	var count int
	err := s.db.QueryRow(query, profileURL).Scan(&count)
	return count > 0, err
}

// UpdateConnectionStatus updates the status of a connection request
func (s *Storage) UpdateConnectionStatus(profileURL, status string) error {
	query := `
		UPDATE connection_requests 
		SET status = ?, accepted_at = CASE WHEN ? = 'accepted' THEN CURRENT_TIMESTAMP ELSE NULL END
		WHERE profile_url = ?
	`
	_, err := s.db.Exec(query, status, status, profileURL)
	return err
}

// GetPendingConnections returns all pending connection requests
func (s *Storage) GetPendingConnections() ([]*models.ConnectionRequest, error) {
	query := `
		SELECT id, profile_url, profile_name, message, sent_at, status
		FROM connection_requests
		WHERE status = 'pending'
		ORDER BY sent_at DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ConnectionRequest
	for rows.Next() {
		req := &models.ConnectionRequest{}
		err := rows.Scan(&req.ID, &req.ProfileURL, &req.ProfileName, &req.Message, &req.SentAt, &req.Status)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// GetAcceptedConnections returns recently accepted connections
func (s *Storage) GetAcceptedConnections(since time.Time) ([]*models.ConnectionRequest, error) {
	query := `
		SELECT id, profile_url, profile_name, message, sent_at, status, accepted_at
		FROM connection_requests
		WHERE status = 'accepted' AND accepted_at >= ?
		ORDER BY accepted_at DESC
	`
	rows, err := s.db.Query(query, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ConnectionRequest
	for rows.Next() {
		req := &models.ConnectionRequest{}
		err := rows.Scan(&req.ID, &req.ProfileURL, &req.ProfileName, &req.Message,
			&req.SentAt, &req.Status, &req.AcceptedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// SaveMessage saves a sent message
func (s *Storage) SaveMessage(msg *models.Message) error {
	query := `
		INSERT INTO messages (profile_url, profile_name, content, sent_at, is_follow_up)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, msg.ProfileURL, msg.ProfileName, msg.Content, msg.SentAt, msg.IsFollowUp)
	return err
}

// MessageExists checks if a message was already sent to a profile
func (s *Storage) MessageExists(profileURL string) (bool, error) {
	query := "SELECT COUNT(*) FROM messages WHERE profile_url = ?"
	var count int
	err := s.db.QueryRow(query, profileURL).Scan(&count)
	return count > 0, err
}

// UpdateActivityStats updates daily activity statistics
func (s *Storage) UpdateActivityStats(date string, field string, increment int) error {
	// Insert or update
	query := fmt.Sprintf(`
		INSERT INTO activity_stats (date, %s) VALUES (?, ?)
		ON CONFLICT(date) DO UPDATE SET %s = %s + ?
	`, field, field, field)
	
	_, err := s.db.Exec(query, date, increment, increment)
	return err
}

// GetActivityStats retrieves activity statistics for a date
func (s *Storage) GetActivityStats(date string) (*models.ActivityStats, error) {
	query := `
		SELECT date, connections_sent, messages_sent, searches_performed, profiles_viewed
		FROM activity_stats
		WHERE date = ?
	`
	stats := &models.ActivityStats{}
	err := s.db.QueryRow(query, date).Scan(
		&stats.Date, &stats.ConnectionsSent, &stats.MessagesSent,
		&stats.SearchesPerformed, &stats.ProfilesViewed,
	)
	if err == sql.ErrNoRows {
		// Return empty stats for date
		return &models.ActivityStats{Date: date}, nil
	}
	return stats, err
}

// GetTodayStats retrieves today's activity statistics
func (s *Storage) GetTodayStats() (*models.ActivityStats, error) {
	today := time.Now().Format("2006-01-02")
	return s.GetActivityStats(today)
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}
