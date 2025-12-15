package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// MouseMover handles human-like mouse movements
type MouseMover struct {
	page *rod.Page
	rand *rand.Rand
}

// NewMouseMover creates a new mouse mover
func NewMouseMover(page *rod.Page) *MouseMover {
	return &MouseMover{
		page: page,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// MoveTo moves the mouse to target coordinates using Bézier curve with human-like behavior
func (m *MouseMover) MoveTo(targetX, targetY float64) error {
	// Get current mouse position (assume starting from a random position if first move)
	currentX := float64(m.rand.Intn(100))
	currentY := float64(m.rand.Intn(100))

	// Generate control points for Bézier curve
	points := m.generateBezierPath(currentX, currentY, targetX, targetY)

	// Move along the path with variable speed
	for i := 0; i < len(points); i++ {
		p := points[i]
		
		// Move mouse to point
		if err := m.page.Mouse.Move(p.X, p.Y, 1); err != nil {
			return err
		}

		// Add micro-delays with variation (simulate human processing time)
		delay := time.Duration(5+m.rand.Intn(10)) * time.Millisecond
		time.Sleep(delay)
	}

	// Add occasional overshoot and correction (human behavior)
	if m.rand.Float64() < 0.3 { // 30% chance of overshoot
		overshootX := targetX + float64(m.rand.Intn(20)-10)
		overshootY := targetY + float64(m.rand.Intn(20)-10)
		m.page.Mouse.Move(overshootX, overshootY, 1)
		time.Sleep(time.Duration(50+m.rand.Intn(100)) * time.Millisecond)
		m.page.Mouse.Move(targetX, targetY, 1)
	}

	return nil
}

// generateBezierPath creates a smooth path using cubic Bézier curves
func (m *MouseMover) generateBezierPath(x1, y1, x4, y4 float64) []Point {
	// Generate random control points for natural curve
	// Control points create the "bend" in the curve
	x2 := x1 + (x4-x1)*m.rand.Float64()
	y2 := y1 + (y4-y1)*m.rand.Float64()
	x3 := x1 + (x4-x1)*m.rand.Float64()
	y3 := y1 + (y4-y1)*m.rand.Float64()

	// Number of points in path (more points = smoother movement)
	steps := int(math.Sqrt(math.Pow(x4-x1, 2)+math.Pow(y4-y1, 2)) / 10)
	if steps < 10 {
		steps = 10
	}
	if steps > 100 {
		steps = 100
	}

	points := make([]Point, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		
		// Cubic Bézier formula: B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
		x := math.Pow(1-t, 3)*x1 +
			3*math.Pow(1-t, 2)*t*x2 +
			3*(1-t)*math.Pow(t, 2)*x3 +
			math.Pow(t, 3)*x4

		y := math.Pow(1-t, 3)*y1 +
			3*math.Pow(1-t, 2)*t*y2 +
			3*(1-t)*math.Pow(t, 2)*y3 +
			math.Pow(t, 3)*y4

		points[i] = Point{X: x, Y: y}
	}

	return points
}

// HoverOverElement simulates random hovering behavior
func (m *MouseMover) HoverOverElement(element *rod.Element) error {
	box, err := element.Shape()
	if err != nil {
		return err
	}

	// Get a random point within the element
	point := box.OnePointInside()
	
	// Move to the element
	if err := m.MoveTo(point.X, point.Y); err != nil {
		return err
	}

	// Hover for a realistic duration
	hoverTime := time.Duration(500+m.rand.Intn(1500)) * time.Millisecond
	time.Sleep(hoverTime)

	return nil
}

// RandomMouseMovement adds random mouse movements to simulate natural behavior
func (m *MouseMover) RandomMouseMovement() error {
	// Get viewport dimensions
	viewport := m.page.MustEval("({width: window.innerWidth, height: window.innerHeight})").Map()
	width := int(viewport["width"].(float64))
	height := int(viewport["height"].(float64))

	// Move to random location
	randomX := float64(m.rand.Intn(width))
	randomY := float64(m.rand.Intn(height))

	return m.MoveTo(randomX, randomY)
}

// ClickElement clicks an element with human-like behavior
func (m *MouseMover) ClickElement(element *rod.Element) error {
	// Hover over element first
	if err := m.HoverOverElement(element); err != nil {
		return err
	}

	// Random delay before clicking (thinking time)
	thinkTime := time.Duration(200+m.rand.Intn(800)) * time.Millisecond
	time.Sleep(thinkTime)

	// Perform click
	if err := element.Click(input.ButtonLeft, 1); err != nil {
		return err
	}

	return nil
}
