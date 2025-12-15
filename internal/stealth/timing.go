package stealth

import (
	"math/rand"
	"time"

	"github.com/keerthana/linkedin-automation-poc/pkg/models"
)

// TimingController manages delays and timing patterns
type TimingController struct {
	config *models.StealthConfig
	rand   *rand.Rand
}

// NewTimingController creates a new timing controller
func NewTimingController(config *models.StealthConfig) *TimingController {
	return &TimingController{
		config: config,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RandomDelay adds a random delay between min and max configured delays
func (t *TimingController) RandomDelay() {
	minMs := t.config.MinActionDelay.Milliseconds()
	maxMs := t.config.MaxActionDelay.Milliseconds()
	
	if maxMs <= minMs {
		time.Sleep(t.config.MinActionDelay)
		return
	}

	delayMs := minMs + t.rand.Int63n(maxMs-minMs)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
}

// ThinkingDelay simulates human thinking/processing time
func (t *TimingController) ThinkingDelay() {
	// 1-4 seconds of "thinking"
	delay := time.Duration(1000+t.rand.Intn(3000)) * time.Millisecond
	time.Sleep(delay)
}

// ReadingDelay simulates reading time based on content length
func (t *TimingController) ReadingDelay(contentLength int) {
	// Average reading speed: 200-250 words per minute
	// Assume ~5 characters per word
	words := contentLength / 5
	
	// Base reading time (milliseconds per word)
	baseTime := 240 + t.rand.Intn(60) // 200-250 WPM range
	
	readingTime := time.Duration(words*baseTime) * time.Millisecond
	
	// Add some randomness (70-130% of calculated time)
	variance := 0.7 + t.rand.Float64()*0.6
	finalTime := time.Duration(float64(readingTime) * variance)
	
	// Cap maximum reading time
	maxReadingTime := 10 * time.Second
	if finalTime > maxReadingTime {
		finalTime = maxReadingTime
	}
	
	time.Sleep(finalTime)
}

// IsBusinessHours checks if current time is within business hours
func (t *TimingController) IsBusinessHours() bool {
	if !t.config.BusinessHoursOnly {
		return true
	}

	now := time.Now()
	currentHour := now.Hour()

	return currentHour >= t.config.BusinessHoursStart && currentHour < t.config.BusinessHoursEnd
}

// WaitForBusinessHours waits until business hours if required
func (t *TimingController) WaitForBusinessHours() {
	if !t.config.BusinessHoursOnly {
		return
	}

	for !t.IsBusinessHours() {
		now := time.Now()
		currentHour := now.Hour()
		
		var hoursToWait int
		if currentHour < t.config.BusinessHoursStart {
			hoursToWait = t.config.BusinessHoursStart - currentHour
		} else {
			// After business hours, wait until next day
			hoursToWait = 24 - currentHour + t.config.BusinessHoursStart
		}
		
		waitDuration := time.Duration(hoursToWait) * time.Hour
		time.Sleep(waitDuration)
	}
}

// RateLimiter manages action rate limits
type RateLimiter struct {
	config           *models.RateLimitsConfig
	dailyConnections int
	dailyMessages    int
	hourlySearches   int
	lastReset        time.Time
	lastHourReset    time.Time
	lastAction       time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *models.RateLimitsConfig) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		config:        config,
		lastReset:     now,
		lastHourReset: now,
		lastAction:    now,
	}
}

// CanSendConnection checks if we can send a connection request
func (r *RateLimiter) CanSendConnection() bool {
	r.resetIfNeeded()
	return r.dailyConnections < r.config.MaxConnectionsPerDay
}

// CanSendMessage checks if we can send a message
func (r *RateLimiter) CanSendMessage() bool {
	r.resetIfNeeded()
	return r.dailyMessages < r.config.MaxMessagesPerDay
}

// CanSearch checks if we can perform a search
func (r *RateLimiter) CanSearch() bool {
	r.resetIfNeeded()
	return r.hourlySearches < r.config.MaxSearchesPerHour
}

// RecordConnection records a connection request
func (r *RateLimiter) RecordConnection() {
	r.resetIfNeeded()
	r.dailyConnections++
	r.lastAction = time.Now()
}

// RecordMessage records a message sent
func (r *RateLimiter) RecordMessage() {
	r.resetIfNeeded()
	r.dailyMessages++
	r.lastAction = time.Now()
}

// RecordSearch records a search performed
func (r *RateLimiter) RecordSearch() {
	r.resetIfNeeded()
	r.hourlySearches++
	r.lastAction = time.Now()
}

// WaitForCooldown waits for the cooldown period after last action
func (r *RateLimiter) WaitForCooldown() {
	elapsed := time.Since(r.lastAction)
	if elapsed < r.config.CooldownPeriod {
		time.Sleep(r.config.CooldownPeriod - elapsed)
	}
}

// resetIfNeeded resets counters if time period has elapsed
func (r *RateLimiter) resetIfNeeded() {
	now := time.Now()
	
	// Reset daily counters at midnight
	if now.Day() != r.lastReset.Day() {
		r.dailyConnections = 0
		r.dailyMessages = 0
		r.lastReset = now
	}
	
	// Reset hourly counters
	if now.Sub(r.lastHourReset) >= time.Hour {
		r.hourlySearches = 0
		r.lastHourReset = now
	}
}

// GetStats returns current rate limit statistics
func (r *RateLimiter) GetStats() map[string]int {
	r.resetIfNeeded()
	return map[string]int{
		"daily_connections": r.dailyConnections,
		"daily_messages":    r.dailyMessages,
		"hourly_searches":   r.hourlySearches,
	}
}
