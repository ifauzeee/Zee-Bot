package handlers

import (
	"testing"
)

func TestToBoldSans(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal Text",
			input:    "Hello",
			expected: "𝗛𝗲𝗹𝗹𝗼",
		},
		{
			name:     "Numbers",
			input:    "123",
			expected: "𝟭𝟮𝟯",
		},
		{
			name:     "Mixed Symbols",
			input:    "A-1",
			expected: "𝗔-𝟭",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBoldSans(tt.input); got != tt.expected {
				t.Errorf("ToBoldSans(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestStyleBuilder_Build(t *testing.T) {
	sb := &StyleBuilder{
		title: "Test",
		icon:  "🧪",
	}

	sb.AddRow("Key", "Value")
	output := sb.Build()

	if len(output) == 0 {
		t.Error("Build() returned empty string")
	}
}
