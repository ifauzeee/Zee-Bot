package utility

import (
	"strings"
	"testing"
)

func TestTimeFormat(t *testing.T) {
}

func TestStringHelpers(t *testing.T) {
	input := "hello world"
	expected := "HELLO WORLD"
	if result := strings.ToUpper(input); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
