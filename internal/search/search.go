package search

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

// Searcher handles LinkedIn profile searches
type Searcher struct {
	page     *rod.Page
	storage  *storage.Storage
	logger   *logrus.Logger
	scroller *stealth.ScrollBehavior
	timing   *stealth.TimingController
	baseURL  string
}

// NewSearcher creates a new searcher
func NewSearcher(page *rod.Page, store *storage.Storage, logger *logrus.Logger, config *models.StealthConfig, baseURL string) *Searcher {
	return &Searcher{
		page:     page,
		storage:  store,
		logger:   logger,
		scroller: stealth.NewScrollBehavior(page),
		timing:   stealth.NewTimingController(config),
		baseURL:  baseURL,
	}
}

// SearchProfiles searches for LinkedIn profiles based on criteria
func (s *Searcher) SearchProfiles(criteria *models.SearchCriteria) ([]*models.Profile, error) {
	s.logger.Infof("Starting profile search with criteria: %+v", criteria)

	// Build search query
	searchQuery := s.buildSearchQuery(criteria)
	searchURL := fmt.Sprintf("%s/search/results/people/?keywords=%s", s.baseURL, searchQuery)

	s.logger.Debugf("Navigating to search URL: %s", searchURL)
	if err := s.page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to search page: %w", err)
	}

	s.page.MustWaitLoad()
	s.timing.RandomDelay()

	// Scroll randomly to simulate reading
	s.scroller.ScrollRandomly()

	// Extract profiles from search results
	profiles, err := s.extractProfilesFromPage(criteria.MaxResults)
	if err != nil {
		return nil, fmt.Errorf("failed to extract profiles: %w", err)
	}

	s.logger.Infof("Found %d profiles matching criteria", len(profiles))
	return profiles, nil
}

// buildSearchQuery builds a search query string from criteria
func (s *Searcher) buildSearchQuery(criteria *models.SearchCriteria) string {
	var parts []string

	if criteria.JobTitle != "" {
		parts = append(parts, criteria.JobTitle)
	}

	if criteria.Company != "" {
		parts = append(parts, criteria.Company)
	}

	if criteria.Location != "" {
		parts = append(parts, criteria.Location)
	}

	for _, keyword := range criteria.Keywords {
		parts = append(parts, keyword)
	}

	return strings.Join(parts, " ")
}

// extractProfilesFromPage extracts profile information from current page
func (s *Searcher) extractProfilesFromPage(maxResults int) ([]*models.Profile, error) {
	var profiles []*models.Profile
	seenURLs := make(map[string]bool)

	// Find all profile cards
	elements, err := s.page.Elements(".entity-result__item")
	if err != nil {
		s.logger.Warn("Could not find profile elements with primary selector, trying alternative")
		elements, err = s.page.Elements(".reusable-search__result-container")
		if err != nil {
			return nil, fmt.Errorf("failed to find profile elements: %w", err)
		}
	}

	for i, element := range elements {
		if len(profiles) >= maxResults {
			break
		}

		// Extract profile information
		profile, err := s.extractProfileInfo(element)
		if err != nil {
			s.logger.Warnf("Failed to extract profile %d: %v", i+1, err)
			continue
		}

		// Check for duplicates
		if seenURLs[profile.ProfileURL] {
			s.logger.Debugf("Skipping duplicate profile: %s", profile.ProfileURL)
			continue
		}

		// Check if profile already exists in database
		exists, err := s.storage.ProfileExists(profile.ProfileURL)
		if err != nil {
			s.logger.Warnf("Error checking profile existence: %v", err)
		} else if exists {
			s.logger.Debugf("Profile already in database: %s", profile.ProfileURL)
			continue
		}

		seenURLs[profile.ProfileURL] = true
		profiles = append(profiles, profile)

		// Save profile to database
		if err := s.storage.SaveProfile(profile); err != nil {
			s.logger.Warnf("Failed to save profile: %v", err)
		}

		// Random delay between processing profiles
		s.timing.RandomDelay()
	}

	return profiles, nil
}

// extractProfileInfo extracts profile details from a search result element
func (s *Searcher) extractProfileInfo(element *rod.Element) (*models.Profile, error) {
	profile := &models.Profile{
		DiscoveredAt: time.Now(),
	}

	// Extract profile URL
	linkElement, err := element.Element("a.app-aware-link")
	if err != nil {
		return nil, fmt.Errorf("failed to find profile link: %w", err)
	}

	href, err := linkElement.Attribute("href")
	if err != nil || href == nil {
		return nil, fmt.Errorf("failed to get profile URL")
	}
	profile.ProfileURL = *href

	// Extract name
	nameElement, err := element.Element(".entity-result__title-text a span")
	if err == nil {
		nameText, _ := nameElement.Text()
		profile.Name = strings.TrimSpace(nameText)
	}

	// Extract title
	titleElement, err := element.Element(".entity-result__primary-subtitle")
	if err == nil {
		titleText, _ := titleElement.Text()
		profile.Title = strings.TrimSpace(titleText)
	}

	// Extract company (often in secondary subtitle)
	companyElement, err := element.Element(".entity-result__secondary-subtitle")
	if err == nil {
		companyText, _ := companyElement.Text()
		profile.Company = strings.TrimSpace(companyText)
	}

	// Extract location
	locationElement, err := element.Element(".entity-result__location")
	if err == nil {
		locationText, _ := locationElement.Text()
		profile.Location = strings.TrimSpace(locationText)
	}

	return profile, nil
}

// HandlePagination navigates through search result pages
func (s *Searcher) HandlePagination(maxPages int) error {
	for page := 1; page < maxPages; page++ {
		// Look for next button
		nextBtn, err := s.page.Timeout(3 * time.Second).Element("button[aria-label='Next']")
		if err != nil {
			s.logger.Info("No more pages available")
			return nil
		}

		// Scroll to next button
		s.scroller.ScrollToElement(nextBtn)
		s.timing.RandomDelay()

		// Click next button
		if err := nextBtn.Click(rod.ButtonLeft, 1); err != nil {
			return fmt.Errorf("failed to click next button: %w", err)
		}

		s.page.MustWaitLoad()
		s.timing.ThinkingDelay()

		// Scroll randomly on new page
		s.scroller.ScrollRandomly()
	}

	return nil
}
