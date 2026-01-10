package mocks

import "snakehem/internal/interfaces"

// MockRandomSource provides deterministic random numbers for testing
type MockRandomSource struct {
	IntValues []int
	IntIndex  int
}

// NewMockRandomSource creates a new mock random source with predefined values
func NewMockRandomSource(values []int) *MockRandomSource {
	return &MockRandomSource{
		IntValues: values,
		IntIndex:  0,
	}
}

// IntN returns the next predefined value (cycles through the list)
func (m *MockRandomSource) IntN(n int) int {
	if len(m.IntValues) == 0 {
		return 0
	}
	val := m.IntValues[m.IntIndex%len(m.IntValues)]
	m.IntIndex++
	return val % n // Ensure it's within bounds
}

// Reset resets the index to start from the beginning
func (m *MockRandomSource) Reset() {
	m.IntIndex = 0
}

// Verify interface compliance at compile time
var _ interfaces.RandomSource = (*MockRandomSource)(nil)
