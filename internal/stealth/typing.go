package stealth

import (
	"math/rand"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// TypingSimulator handles human-like typing behavior
type TypingSimulator struct {
	page *rod.Page
	rand *rand.Rand
}

// NewTypingSimulator creates a new typing simulator
func NewTypingSimulator(page *rod.Page) *TypingSimulator {
	return &TypingSimulator{
		page: page,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TypeText types text with human-like characteristics
func (t *TypingSimulator) TypeText(element *rod.Element, text string) error {
	// Focus on the element first
	if err := element.Focus(); err != nil {
		return err
	}

	// Add initial thinking delay
	time.Sleep(time.Duration(300+t.rand.Intn(700)) * time.Millisecond)

	// Type each character with realistic timing
	for i, char := range text {
		// Simulate occasional typos (5% chance)
		if t.rand.Float64() < 0.05 && i < len(text)-1 {
			// Type wrong character
			wrongChar := t.getRandomChar()
			element.MustInput(string(wrongChar))
			
			// Pause (realize mistake)
			time.Sleep(time.Duration(200+t.rand.Intn(400)) * time.Millisecond)
			
			// Backspace to correct
			t.page.Keyboard.MustPress(rod.Backspace)
			time.Sleep(time.Duration(100+t.rand.Intn(200)) * time.Millisecond)
		}

		// Type the actual character
		element.MustInput(string(char))

		// Variable delay between keystrokes
		delay := t.getKeystrokeDelay(char, i, len(text))
		time.Sleep(delay)
	}

	return nil
}

// getKeystrokeDelay calculates realistic delay between keystrokes
func (t *TypingSimulator) getKeystrokeDelay(char rune, position, textLength int) time.Duration {
	// Base delay: 80-150ms (average human typing speed)
	baseDelay := 80 + t.rand.Intn(70)

	// Add extra delay for specific characters
	switch char {
	case ' ':
		// Spaces often have slightly longer delays
		baseDelay += t.rand.Intn(30)
	case '.', ',', '!', '?':
		// Punctuation often has pauses
		baseDelay += t.rand.Intn(100)
	case '\n':
		// Newlines have longer delays
		baseDelay += 200 + t.rand.Intn(300)
	}

	// First few characters are slower (thinking/starting)
	if position < 5 {
		baseDelay += t.rand.Intn(50)
	}

	// Occasional longer pauses (thinking/rereading)
	if t.rand.Float64() < 0.1 {
		baseDelay += 200 + t.rand.Intn(500)
	}

	return time.Duration(baseDelay) * time.Millisecond
}

// getRandomChar returns a random character for typo simulation
func (t *TypingSimulator) getRandomChar() rune {
	chars := "abcdefghijklmnopqrstuvwxyz"
	return rune(chars[t.rand.Intn(len(chars))])
}

// TypeWithBackspace simulates typing with occasional backspacing
func (t *TypingSimulator) TypeWithBackspace(element *rod.Element, text string) error {
	words := strings.Fields(text)
	
	for i, word := range words {
		// Type the word
		if err := t.TypeText(element, word); err != nil {
			return err
		}

		// 15% chance to backspace and retype part of word
		if t.rand.Float64() < 0.15 && len(word) > 3 {
			// Backspace a few characters
			backspaceCount := 1 + t.rand.Intn(3)
			for j := 0; j < backspaceCount; j++ {
				t.page.Keyboard.MustPress(rod.Backspace)
				time.Sleep(time.Duration(80+t.rand.Intn(120)) * time.Millisecond)
			}

			// Retype the characters
			retypeChars := word[len(word)-backspaceCount:]
			t.TypeText(element, retypeChars)
		}

		// Add space between words (except last word)
		if i < len(words)-1 {
			element.MustInput(" ")
			time.Sleep(time.Duration(100+t.rand.Intn(150)) * time.Millisecond)
		}
	}

	return nil
}

// FillForm fills a form field with human-like typing
func (t *TypingSimulator) FillForm(selector string, text string) error {
	element, err := t.page.Element(selector)
	if err != nil {
		return err
	}

	// Clear existing text first
	element.MustSelectAllText()
	time.Sleep(time.Duration(50+t.rand.Intn(100)) * time.Millisecond)
	t.page.Keyboard.MustPress(rod.Backspace)
	
	// Type new text
	return t.TypeText(element, text)
}
