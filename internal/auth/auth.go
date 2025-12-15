package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/keerthana/linkedin-automation-poc/internal/stealth"
	"github.com/keerthana/linkedin-automation-poc/pkg/models"
	"github.com/sirupsen/logrus"
)

// Authenticator handles LinkedIn authentication
type Authenticator struct {
	config  *models.LinkedInConfig
	page    *rod.Page
	logger  *logrus.Logger
	typer   *stealth.TypingSimulator
	scroller *stealth.ScrollBehavior
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(config *models.LinkedInConfig, page *rod.Page, logger *logrus.Logger) *Authenticator {
	return &Authenticator{
		config:   config,
		page:     page,
		logger:   logger,
		typer:    stealth.NewTypingSimulator(page),
		scroller: stealth.NewScrollBehavior(page),
	}
}

// Login performs LinkedIn login with credentials
func (a *Authenticator) Login(email, password string) error {
	a.logger.Info("Starting LinkedIn login process")

	// Check if we can restore session
	if err := a.RestoreSession(); err == nil {
		a.logger.Info("Session restored successfully")
		if a.IsLoggedIn() {
			return nil
		}
		a.logger.Warn("Restored session is invalid, performing fresh login")
	}

	// Navigate to LinkedIn login page
	a.logger.Debug("Navigating to LinkedIn login page")
	if err := a.page.Navigate(a.config.BaseURL + "/login"); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Wait for page to load
	a.page.MustWaitLoad()
	time.Sleep(2 * time.Second)

	// Find and fill email field
	a.logger.Debug("Filling email field")
	emailField, err := a.page.Element("#username")
	if err != nil {
		return fmt.Errorf("failed to find email field: %w", err)
	}

	if err := a.typer.TypeText(emailField, email); err != nil {
		return fmt.Errorf("failed to type email: %w", err)
	}

	// Random delay between fields
	time.Sleep(time.Duration(500+time.Now().UnixNano()%1000) * time.Millisecond)

	// Find and fill password field
	a.logger.Debug("Filling password field")
	passwordField, err := a.page.Element("#password")
	if err != nil {
		return fmt.Errorf("failed to find password field: %w", err)
	}

	if err := a.typer.TypeText(passwordField, password); err != nil {
		return fmt.Errorf("failed to type password: %w", err)
	}

	// Random delay before clicking submit
	time.Sleep(time.Duration(800+time.Now().UnixNano()%1200) * time.Millisecond)

	// Click sign in button
	a.logger.Debug("Clicking sign in button")
	signInBtn, err := a.page.Element("button[type='submit']")
	if err != nil {
		return fmt.Errorf("failed to find sign in button: %w", err)
	}

	if err := signInBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click sign in button: %w", err)
	}

	// Wait for navigation
	time.Sleep(5 * time.Second)

	// Check for security challenges
	if err := a.handleSecurityChallenges(); err != nil {
		return fmt.Errorf("security challenge detected: %w", err)
	}

	// Verify login success
	if !a.IsLoggedIn() {
		return fmt.Errorf("login failed: not authenticated after login attempt")
	}

	a.logger.Info("Login successful")

	// Save session for future use
	if err := a.SaveSession(); err != nil {
		a.logger.Warnf("Failed to save session: %v", err)
	}

	return nil
}

// IsLoggedIn checks if the user is currently logged in
func (a *Authenticator) IsLoggedIn() bool {
	// Check for presence of navigation elements that only appear when logged in
	_, err := a.page.Timeout(3 * time.Second).Element("nav.global-nav")
	if err != nil {
		// Try alternative selector
		_, err = a.page.Timeout(3 * time.Second).Element("#global-nav")
		return err == nil
	}
	return true
}

// handleSecurityChallenges detects and handles security challenges
func (a *Authenticator) handleSecurityChallenges() error {
	// Check for 2FA/verification code prompt
	verificationElement, err := a.page.Timeout(3 * time.Second).Element("#input__email_verification_pin")
	if err == nil && verificationElement != nil {
		a.logger.Error("Two-factor authentication detected - manual intervention required")
		return fmt.Errorf("2FA verification required - please complete manually")
	}

	// Check for CAPTCHA
	captchaElement, err := a.page.Timeout(3 * time.Second).Element("#captcha-internal")
	if err == nil && captchaElement != nil {
		a.logger.Error("CAPTCHA detected - manual intervention required")
		return fmt.Errorf("CAPTCHA verification required - please complete manually")
	}

	// Check for security checkpoint
	checkpointElement, err := a.page.Timeout(3 * time.Second).Element("[data-test-id='security-challenge']")
	if err == nil && checkpointElement != nil {
		a.logger.Error("Security checkpoint detected - manual intervention required")
		return fmt.Errorf("security checkpoint detected - please complete manually")
	}

	return nil
}

// SaveSession saves cookies to file for session persistence
func (a *Authenticator) SaveSession() error {
	cookies := a.page.MustCookies()

	// Create session directory if it doesn't exist
	sessionDir := a.config.SessionPath
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionFile := filepath.Join(sessionDir, "linkedin_session.json")

	// Marshal cookies to JSON
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	// Write to file
	if err := os.WriteFile(sessionFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	a.logger.Debug("Session saved successfully")
	return nil
}

// RestoreSession restores cookies from saved session
func (a *Authenticator) RestoreSession() error {
	sessionFile := filepath.Join(a.config.SessionPath, "linkedin_session.json")

	// Check if session file exists
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		return fmt.Errorf("no saved session found")
	}

	// Read session file
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	// Unmarshal cookies
	var cookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("failed to unmarshal cookies: %w", err)
	}

	// Set cookies
	if err := a.page.SetCookies(cookies); err != nil {
		return fmt.Errorf("failed to set cookies: %w", err)
	}

	a.logger.Debug("Session restored from file")
	return nil
}

// Logout performs logout from LinkedIn
func (a *Authenticator) Logout() error {
	a.logger.Info("Logging out from LinkedIn")

	// Navigate to logout
	if err := a.page.Navigate(a.config.BaseURL + "/m/logout"); err != nil {
		return fmt.Errorf("failed to navigate to logout: %w", err)
	}

	time.Sleep(2 * time.Second)

	// Clear saved session
	sessionFile := filepath.Join(a.config.SessionPath, "linkedin_session.json")
	os.Remove(sessionFile)

	return nil
}
