package messaging

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/keerthana/linkedin-automation-poc/internal/stealth"
	"github.com/keerthana/linkedin-automation-poc/internal/storage"
	"github.com/keerthana/linkedin-automation-poc/pkg/models"
	"github.com/sirupsen/logrus"
)

// Messenger handles connection requests and messaging
type Messenger struct {
	page       *rod.Page
	storage    *storage.Storage
	logger     *logrus.Logger
	typer      *stealth.TypingSimulator
	mover      *stealth.MouseMover
	scroller   *stealth.ScrollBehavior
	timing     *stealth.TimingController
	rateLimiter *stealth.RateLimiter
}

// NewMessenger creates a new messenger
func NewMessenger(page *rod.Page, store *storage.Storage, logger *logrus.Logger, 
	stealthConfig *models.StealthConfig, rateLimitsConfig *models.RateLimitsConfig) *Messenger {
	return &Messenger{
		page:        page,
		storage:     store,
		logger:      logger,
		typer:       stealth.NewTypingSimulator(page),
		mover:       stealth.NewMouseMover(page),
		scroller:    stealth.NewScrollBehavior(page),
		timing:      stealth.NewTimingController(stealthConfig),
		rateLimiter: stealth.NewRateLimiter(rateLimitsConfig),
	}
}

// SendConnectionRequest sends a connection request to a profile
func (m *Messenger) SendConnectionRequest(profileURL, profileName, message string) error {
	// Check rate limits
	if !m.rateLimiter.CanSendConnection() {
		return fmt.Errorf("daily connection limit reached")
	}

	// Check if already sent
	exists, err := m.storage.ConnectionRequestExists(profileURL)
	if err != nil {
		m.logger.Warnf("Error checking connection request: %v", err)
	} else if exists {
		return fmt.Errorf("connection request already sent to this profile")
	}

	m.logger.Infof("Sending connection request to: %s", profileName)

	// Navigate to profile
	if err := m.page.Navigate(profileURL); err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	m.page.MustWaitLoad()
	m.timing.RandomDelay()

	// Scroll to simulate reading profile
	m.scroller.ScrollRandomly()
	m.timing.ReadingDelay(500) // Simulate reading profile

	// Find Connect button
	connectBtn, err := m.findConnectButton()
	if err != nil {
		return fmt.Errorf("failed to find connect button: %w", err)
	}

	// Scroll to button
	m.scroller.ScrollToElement(connectBtn)
	m.timing.ThinkingDelay()

	// Click Connect button using human-like mouse movement
	if err := m.mover.ClickElement(connectBtn); err != nil {
		return fmt.Errorf("failed to click connect button: %w", err)
	}

	m.timing.RandomDelay()

	// Check if "Add a note" option appears
	addNoteBtn, err := m.page.Timeout(3 * time.Second).Element("button[aria-label='Add a note']")
	if err == nil && message != "" {
		// Click "Add a note"
		if err := m.mover.ClickElement(addNoteBtn); err != nil {
			m.logger.Warnf("Failed to click add note button: %v", err)
		} else {
			m.timing.RandomDelay()

			// Find message textarea
			messageField, err := m.page.Element("#custom-message")
			if err == nil {
				// Type personalized message
				if err := m.typer.TypeText(messageField, message); err != nil {
					m.logger.Warnf("Failed to type message: %v", err)
				}
				m.timing.ThinkingDelay()
			}
		}
	}

	// Click Send button
	sendBtn, err := m.page.Timeout(3 * time.Second).Element("button[aria-label='Send now']")
	if err != nil {
		// Try alternative selector
		sendBtn, err = m.page.Element("button[aria-label='Send']")
		if err != nil {
			return fmt.Errorf("failed to find send button: %w", err)
		}
	}

	if err := m.mover.ClickElement(sendBtn); err != nil {
		return fmt.Errorf("failed to click send button: %w", err)
	}

	m.logger.Info("Connection request sent successfully")

	// Save to database
	req := &models.ConnectionRequest{
		ProfileURL:  profileURL,
		ProfileName: profileName,
		Message:     message,
		SentAt:      time.Now(),
		Status:      "pending",
	}

	if err := m.storage.SaveConnectionRequest(req); err != nil {
		m.logger.Warnf("Failed to save connection request: %v", err)
	}

	// Record for rate limiting
	m.rateLimiter.RecordConnection()

	// Update activity stats
	today := time.Now().Format("2006-01-02")
	m.storage.UpdateActivityStats(today, "connections_sent", 1)

	return nil
}

// findConnectButton finds the connect button on a profile page
func (m *Messenger) findConnectButton() (*rod.Element, error) {
	// Try various selectors for Connect button
	selectors := []string{
		"button[aria-label*='Connect']",
		"button.pvs-profile-actions__action[aria-label*='Connect']",
		"button:has-text('Connect')",
	}

	for _, selector := range selectors {
		btn, err := m.page.Timeout(3 * time.Second).Element(selector)
		if err == nil {
			return btn, nil
		}
	}

	return nil, fmt.Errorf("connect button not found")
}

// SendMessage sends a message to a connected profile
func (m *Messenger) SendMessage(profileURL, profileName, message string) error {
	// Check rate limits
	if !m.rateLimiter.CanSendMessage() {
		return fmt.Errorf("daily message limit reached")
	}

	m.logger.Infof("Sending message to: %s", profileName)

	// Navigate to profile
	if err := m.page.Navigate(profileURL); err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	m.page.MustWaitLoad()
	m.timing.RandomDelay()

	// Find Message button
	messageBtn, err := m.findMessageButton()
	if err != nil {
		return fmt.Errorf("failed to find message button: %w", err)
	}

	// Click Message button
	if err := m.mover.ClickElement(messageBtn); err != nil {
		return fmt.Errorf("failed to click message button: %w", err)
	}

	m.timing.RandomDelay()

	// Find message input field
	messageField, err := m.page.Timeout(5 * time.Second).Element(".msg-form__contenteditable")
	if err != nil {
		return fmt.Errorf("failed to find message field: %w", err)
	}

	// Type message with human-like behavior
	if err := m.typer.TypeWithBackspace(messageField, message); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	m.timing.ThinkingDelay()

	// Find and click send button
	sendBtn, err := m.page.Element("button[type='submit'].msg-form__send-button")
	if err != nil {
		return fmt.Errorf("failed to find send button: %w", err)
	}

	if err := m.mover.ClickElement(sendBtn); err != nil {
		return fmt.Errorf("failed to click send button: %w", err)
	}

	m.logger.Info("Message sent successfully")

	// Save to database
	msg := &models.Message{
		ProfileURL:  profileURL,
		ProfileName: profileName,
		Content:     message,
		SentAt:      time.Now(),
		IsFollowUp:  true,
	}

	if err := m.storage.SaveMessage(msg); err != nil {
		m.logger.Warnf("Failed to save message: %v", err)
	}

	// Record for rate limiting
	m.rateLimiter.RecordMessage()

	// Update activity stats
	today := time.Now().Format("2006-01-02")
	m.storage.UpdateActivityStats(today, "messages_sent", 1)

	return nil
}

// findMessageButton finds the message button on a profile page
func (m *Messenger) findMessageButton() (*rod.Element, error) {
	selectors := []string{
		"button[aria-label*='Message']",
		"button.pvs-profile-actions__action[aria-label*='Message']",
		"button:has-text('Message')",
	}

	for _, selector := range selectors {
		btn, err := m.page.Timeout(3 * time.Second).Element(selector)
		if err == nil {
			return btn, nil
		}
	}

	return nil, fmt.Errorf("message button not found")
}

// ProcessTemplate replaces variables in message template
func (m *Messenger) ProcessTemplate(template string, profile *models.Profile) string {
	message := template

	// Replace {{name}} with profile name
	if profile.Name != "" {
		firstName := strings.Split(profile.Name, " ")[0]
		message = strings.ReplaceAll(message, "{{name}}", firstName)
		message = strings.ReplaceAll(message, "{{full_name}}", profile.Name)
	}

	// Replace {{title}} with profile title
	if profile.Title != "" {
		message = strings.ReplaceAll(message, "{{title}}", profile.Title)
	}

	// Replace {{company}} with profile company
	if profile.Company != "" {
		message = strings.ReplaceAll(message, "{{company}}", profile.Company)
	}

	return message
}

// GetRateLimitStats returns current rate limit statistics
func (m *Messenger) GetRateLimitStats() map[string]int {
	return m.rateLimiter.GetStats()
}
