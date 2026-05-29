package utils

import "testing"

func TestToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		expected  int
		wantError bool
	}{
		{"positive number", "42", 42, false},
		{"zero", "0", 0, false},
		{"negative number", "-5", -5, false},
		{"invalid string", "abc", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToInt(tt.input)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %d; want %d", result, tt.expected)
			}
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string input", "hello", "hello"},
		{"number input", 42, "%!s(int=42)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("got %q; want %q", result, tt.expected)
			}
		})
	}
}
