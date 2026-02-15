package tui

import "testing"

func TestClampInt(t *testing.T) {
	tests := []struct {
		name     string
		v        int
		min      int
		max      int
		expected int
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -5, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 10},
		{"negative range within", -5, -10, 0, -5},
		{"negative range below", -15, -10, 0, -10},
		{"negative range above", 5, -10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clampInt(tt.v, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("clampInt(%d, %d, %d) = %d, want %d", tt.v, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}
