package stealth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

// BrowserConfig holds browser stealth configuration
type BrowserConfig struct {
	UserAgent      string
	ViewportWidth  int
	ViewportHeight int
	Headless       bool
}

// SetupStealthBrowser creates a browser with anti-detection measures
func SetupStealthBrowser(config BrowserConfig) (*rod.Browser, error) {
	// Launch browser with stealth settings
	l := launcher.New().
		Headless(config.Headless).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars", "true").
		Set("disable-background-networking", "true").
		Set("disable-background-timer-throttling", "true").
		Set("disable-backgrounding-occluded-windows", "true").
		Set("disable-breakpad", "true").
		Set("disable-client-side-phishing-detection", "true").
		Set("disable-default-apps", "true").
		Set("disable-dev-shm-usage", "true").
		Set("disable-extensions", "true").
		Set("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees").
		Set("disable-hang-monitor", "true").
		Set("disable-ipc-flooding-protection", "true").
		Set("disable-popup-blocking", "true").
		Set("disable-prompt-on-repost", "true").
		Set("disable-renderer-backgrounding", "true").
		Set("disable-sync", "true").
		Set("force-color-profile", "srgb").
		Set("metrics-recording-only", "true").
		Set("no-first-run", "true").
		Set("enable-automation", "false").
		Set("password-store", "basic").
		Set("use-mock-keychain", "true")

	if config.UserAgent != "" {
		l.Set("user-agent", config.UserAgent)
	}

	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()

	return browser, nil
}

// ApplyStealthScripts applies stealth JavaScript to the page
func ApplyStealthScripts(page *rod.Page) error {
	// Use go-rod/stealth package for advanced evasion
	if err := stealth.Apply(page); err != nil {
		return err
	}

	// Additional custom stealth scripts
	_, err := page.Eval(`() => {
		// Remove webdriver property
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});

		// Override plugins
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});

		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});

		// Override chrome property
		window.chrome = {
			runtime: {}
		};

		// Override permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
	}`)

	return err
}

// ScrollBehavior handles human-like scrolling
type ScrollBehavior struct {
	page *rod.Page
	rand *rand.Rand
}

// NewScrollBehavior creates a new scroll behavior handler
func NewScrollBehavior(page *rod.Page) *ScrollBehavior {
	return &ScrollBehavior{
		page: page,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ScrollRandomly performs random scrolling to simulate reading
func (s *ScrollBehavior) ScrollRandomly() error {
	// Get page height
	pageHeight := s.page.MustEval("document.body.scrollHeight").Int()
	viewportHeight := s.page.MustEval("window.innerHeight").Int()

	// Random number of scroll actions
	scrollCount := 3 + s.rand.Intn(7)

	for i := 0; i < scrollCount; i++ {
		// Random scroll distance (simulate reading sections)
		scrollDistance := viewportHeight/2 + s.rand.Intn(viewportHeight/2)
		
		// Get current scroll position
		currentScroll := s.page.MustEval("window.pageYOffset").Int()
		targetScroll := currentScroll + scrollDistance

		// Don't scroll past page end
		if targetScroll > pageHeight-viewportHeight {
			targetScroll = pageHeight - viewportHeight
		}

		// Smooth scroll with variable speed
		s.smoothScrollTo(targetScroll)

		// Pause to "read" content
		readTime := time.Duration(1000+s.rand.Intn(3000)) * time.Millisecond
		time.Sleep(readTime)

		// Occasionally scroll back up (re-reading)
		if s.rand.Float64() < 0.2 {
			backScroll := currentScroll - s.rand.Intn(200)
			if backScroll < 0 {
				backScroll = 0
			}
			s.smoothScrollTo(backScroll)
			time.Sleep(time.Duration(500+s.rand.Intn(1000)) * time.Millisecond)
		}
	}

	return nil
}

// smoothScrollTo performs smooth scrolling to target position
func (s *ScrollBehavior) smoothScrollTo(target int) {
	current := s.page.MustEval("window.pageYOffset").Int()
	distance := target - current
	steps := 20 + s.rand.Intn(10)

	for i := 0; i <= steps; i++ {
		// Ease-out animation curve
		progress := float64(i) / float64(steps)
		easeProgress := 1 - (1-progress)*(1-progress)
		
		scrollPos := current + int(float64(distance)*easeProgress)
		s.page.MustEval(fmt.Sprintf("window.scrollTo(0, %d)", scrollPos))

		// Variable delay between steps
		delay := time.Duration(10+s.rand.Intn(20)) * time.Millisecond
		time.Sleep(delay)
	}
}

// ScrollToElement scrolls an element into view naturally
func (s *ScrollBehavior) ScrollToElement(element *rod.Element) error {
	// Get element position
	box, err := element.Shape()
	if err != nil {
		return err
	}

	// Scroll to slightly above the element (natural reading position)
	targetY := int(box.Box().Y) - 100 - s.rand.Intn(100)
	if targetY < 0 {
		targetY = 0
	}

	s.smoothScrollTo(targetY)

	// Pause to "notice" the element
	time.Sleep(time.Duration(300+s.rand.Intn(500)) * time.Millisecond)

	return nil
}

// RandomUserAgent generates a realistic user agent string
func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}
