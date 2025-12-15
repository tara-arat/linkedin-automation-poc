# LinkedIn Automation Proof-of-Concept

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)

> **⚠️ CRITICAL DISCLAIMER - READ BEFORE PROCEEDING**
> 
> This project is created **exclusively for educational and technical evaluation purposes**. 
> 
> **Using this tool violates LinkedIn's Terms of Service** and may result in:
> - Permanent account suspension
> - Legal action from LinkedIn
> - Loss of professional network and connections
> 
> **DO NOT use this tool on live LinkedIn accounts or in production environments.**
> This is a technical demonstration only.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Anti-Detection Techniques](#anti-detection-techniques)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Technical Implementation](#technical-implementation)
- [Stealth Techniques Explained](#stealth-techniques-explained)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This proof-of-concept demonstrates advanced browser automation capabilities using Go and the Rod library. It showcases sophisticated anti-detection techniques, human-like behavior simulation, and clean architectural patterns for building automation tools.

### Purpose

- **Technical Demonstration**: Showcase browser automation and anti-bot detection techniques
- **Educational Resource**: Learn about stealth automation, fingerprint masking, and behavior simulation
- **Architecture Example**: Clean, modular Go code following best practices

### What This Tool Does

- Authenticates with LinkedIn using session persistence
- Searches for profiles based on criteria (job title, company, location)
- Sends connection requests with personalized messages
- Manages messaging with accepted connections
- Tracks all activities in a local SQLite database
- Implements 8+ anti-detection techniques

## ✨ Features

### Core Functionality

- **🔐 Authentication System**
  - Login with credentials
  - Session cookie persistence
  - Security checkpoint detection (2FA, CAPTCHA)
  - Automatic session restoration

- **🔍 Search & Targeting**
  - Search by job title, company, location, keywords
  - Profile URL collection
  - Pagination handling
  - Duplicate detection

- **🤝 Connection Management**
  - Send connection requests
  - Personalized message templates
  - Daily limit enforcement
  - Request tracking

- **💬 Messaging System**
  - Automated follow-up messages
  - Template variable support (`{{name}}`, `{{company}}`, etc.)
  - Message history tracking

- **💾 State Persistence**
  - SQLite database for all data
  - Track sent requests and messages
  - Activity statistics
  - Resume after interruption

## 🕵️ Anti-Detection Techniques

This project implements **8+ sophisticated stealth techniques** to simulate authentic human behavior:

### 1. **Human-like Mouse Movement** ✅ (Mandatory)
- Bézier curve-based trajectories
- Variable speed and acceleration
- Natural overshoot with micro-corrections
- Randomized path generation

### 2. **Randomized Timing Patterns** ✅ (Mandatory)
- Variable delays between actions (2-5 seconds configurable)
- Realistic thinking time simulation
- Reading time based on content length
- Random pauses during workflows

### 3. **Browser Fingerprint Masking** ✅ (Mandatory)
- Custom user agent strings
- Disabled automation flags (`navigator.webdriver`)
- Modified browser properties
- Plugin and language spoofing
- Chrome runtime emulation

### 4. **Random Scrolling Behavior** ✅
- Variable scroll speeds
- Natural acceleration/deceleration
- Occasional scroll-back movements
- Smooth scrolling with easing curves
- Reading pause simulation

### 5. **Realistic Typing Simulation** ✅
- Variable keystroke intervals (80-150ms)
- Occasional typos with corrections
- Backspace patterns
- Character-specific delays (punctuation, spaces)
- Initial slowness (thinking/starting)

### 6. **Mouse Hovering & Movement** ✅
- Random hover events over elements
- Natural cursor wandering
- Realistic hover duration (500-2000ms)
- Pre-click hover simulation

### 7. **Activity Scheduling** ✅
- Business hours enforcement (9 AM - 5 PM configurable)
- Automatic waiting outside hours
- Timezone-aware scheduling
- Work pattern simulation

### 8. **Rate Limiting & Throttling** ✅
- Connection request quotas (20/day default)
- Message limits (15/day default)
- Search throttling (5/hour default)
- Cooldown periods (30 minutes default)
- Daily/hourly counter resets

## 🏗️ Architecture

### Clean Modular Design

```
linkedin-automation-poc/
├── cmd/
│   └── linkedin-automation/    # Main application entry point
│       └── main.go
├── internal/
│   ├── auth/                   # Authentication & session management
│   ├── config/                 # Configuration loading & validation
│   ├── logger/                 # Structured logging
│   ├── messaging/              # Connections & messaging
│   ├── search/                 # Profile search & extraction
│   ├── stealth/                # Anti-detection techniques
│   │   ├── browser.go          # Fingerprint masking
│   │   ├── mouse.go            # Bézier mouse movement
│   │   ├── typing.go           # Realistic typing
│   │   └── timing.go           # Rate limiting & delays
│   └── storage/                # SQLite persistence
└── pkg/
    └── models/                 # Shared data models
```

### Design Principles

- **Separation of Concerns**: Each package has a single responsibility
- **Dependency Injection**: Components receive dependencies explicitly
- **Error Handling**: Comprehensive error wrapping and logging
- **Configurability**: YAML configuration with environment overrides
- **Testability**: Interfaces for easy mocking and testing

## 📦 Prerequisites

### Required Software

- **Go 1.21 or higher**: [Download Go](https://golang.org/dl/)
- **Google Chrome**: Required for Rod browser automation
- **Git**: For cloning the repository

### Operating System Support

- ✅ Windows 10/11
- ✅ macOS 10.15+
- ✅ Linux (Ubuntu 20.04+, Debian, etc.)

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/keerthana/linkedin-automation-poc.git
cd linkedin-automation-poc
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Configure Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your credentials (USE TEST ACCOUNT ONLY!)
# NEVER use your real LinkedIn account
```

### 4. Build the Application

```bash
# Build for your current platform
go build -o linkedin-automation ./cmd/linkedin-automation

# Or use the Makefile
make build
```

## ⚙️ Configuration

### Environment Variables (.env)

```env
# LinkedIn Credentials (USE TEST ACCOUNT ONLY!)
LINKEDIN_EMAIL=test-account@example.com
LINKEDIN_PASSWORD=your-test-password

# Optional Settings
LINKEDIN_SESSION_PATH=./sessions
DATABASE_PATH=./data/automation.db
LOG_LEVEL=info
```

### Configuration File (config.yaml)

```yaml
linkedin:
  email: ""  # Set via LINKEDIN_EMAIL
  password: ""  # Set via LINKEDIN_PASSWORD
  session_path: "./sessions"
  base_url: "https://www.linkedin.com"

stealth:
  enable_mouse_movement: true
  enable_typing_simulation: true
  enable_random_scrolling: true
  enable_hovering: true
  min_action_delay: 2s
  max_action_delay: 5s
  business_hours_only: true
  business_hours_start: 9   # 9 AM
  business_hours_end: 17    # 5 PM

rate_limits:
  max_connections_per_day: 20
  max_messages_per_day: 15
  max_searches_per_hour: 5
  cooldown_period: 30m

storage:
  database_path: "./data/automation.db"
  backup_path: "./data/backups"

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json or text
  output_path: "./logs/automation.log"
```

## 📖 Usage

### Basic Usage

```bash
# Run with default configuration
./linkedin-automation

# Run with custom config file
./linkedin-automation -config=custom-config.yaml
```

### Example Workflow

The main.go file includes an example workflow that:
1. Authenticates with LinkedIn
2. Searches for profiles matching criteria
3. Sends personalized connection requests
4. Tracks all activities in the database

### Customizing the Workflow

Edit `cmd/linkedin-automation/main.go` to customize:

```go
// Modify search criteria
searchCriteria := &models.SearchCriteria{
    JobTitle:   "Product Manager",
    Company:    "Google",
    Location:   "New York",
    MaxResults: 20,
}

// Customize connection message
connectionMessage := "Hi {{name}}, I noticed you work at {{company}} as a {{title}}. Would love to connect!"
```

## 📁 Project Structure

```
.
├── cmd/
│   └── linkedin-automation/
│       └── main.go                 # Application entry point
├── internal/
│   ├── auth/
│   │   └── auth.go                 # Authentication logic
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── logger/
│   │   └── logger.go               # Logging setup
│   ├── messaging/
│   │   └── messaging.go            # Connection & messaging
│   ├── search/
│   │   └── search.go               # Profile search
│   ├── stealth/
│   │   ├── browser.go              # Browser stealth
│   │   ├── mouse.go                # Mouse movement
│   │   ├── typing.go               # Typing simulation
│   │   └── timing.go               # Timing & rate limits
│   └── storage/
│       └── storage.go              # Database operations
├── pkg/
│   └── models/
│       ├── config.go               # Configuration models
│       └── models.go               # Domain models
├── .env.example                    # Example environment file
├── .gitignore                      # Git ignore rules
├── config.yaml                     # Configuration file
├── go.mod                          # Go module definition
├── Makefile                        # Build automation
└── README.md                       # This file
```

## 🔧 Technical Implementation

### Stealth Techniques Explained

#### 1. Bézier Curve Mouse Movement

Uses cubic Bézier curves to create natural mouse trajectories:

```go
// Formula: B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
x := math.Pow(1-t, 3)*x1 + 3*math.Pow(1-t, 2)*t*x2 + 
     3*(1-t)*math.Pow(t, 2)*x3 + math.Pow(t, 3)*x4
```

#### 2. Browser Fingerprint Masking

Modifies browser properties to avoid detection:

```javascript
Object.defineProperty(navigator, 'webdriver', {
  get: () => undefined
});
```

#### 3. Realistic Typing

Simulates human typing with:
- Variable keystroke delays (80-150ms)
- Occasional typos and corrections
- Character-specific pauses
- Initial slowness

### Database Schema

```sql
CREATE TABLE profiles (
    id INTEGER PRIMARY KEY,
    profile_url TEXT UNIQUE,
    name TEXT,
    title TEXT,
    company TEXT,
    location TEXT,
    discovered_at DATETIME
);

CREATE TABLE connection_requests (
    id INTEGER PRIMARY KEY,
    profile_url TEXT UNIQUE,
    profile_name TEXT,
    message TEXT,
    sent_at DATETIME,
    status TEXT,
    accepted_at DATETIME
);
```

## 🐛 Troubleshooting

### Common Issues

**Issue**: "Failed to launch browser"
```
Solution: Ensure Google Chrome is installed and accessible
```

**Issue**: "Authentication failed"
```
Solution: 
1. Check your credentials in .env
2. Verify LinkedIn isn't blocking login attempts
3. Complete any security challenges manually
```

**Issue**: "Rate limit reached"
```
Solution: This is intentional! Adjust limits in config.yaml
```

### Debug Mode

Enable debug logging:

```yaml
logging:
  level: "debug"
```

## 🤝 Contributing

This is an educational project. Contributions that enhance the technical demonstration are welcome:

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Rod](https://github.com/go-rod/rod) - Excellent browser automation library
- [Rod Stealth](https://github.com/go-rod/stealth) - Anti-detection evasion techniques

## ⚖️ Legal & Ethical Notice

**IMPORTANT**: This tool is provided for educational purposes to demonstrate:
- Browser automation techniques
- Anti-detection mechanisms
- Clean Go architecture patterns

**You are solely responsible** for how you use this software. The author assumes no liability for misuse, account bans, legal issues, or any other consequences.

**Recommended Use**:
- Study the code to learn automation techniques
- Use as a reference for building legitimate automation tools
- Create demonstration videos for educational content
- DO NOT use on real LinkedIn accounts

---

**Built with ❤️ for educational purposes only**
