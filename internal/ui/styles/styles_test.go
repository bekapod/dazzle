package styles_test

import (
	"testing"

	"dazzle/internal/ui/styles"
)

func TestMethodColor_KnownMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range methods {
		color := styles.MethodColor(m)
		// Each known method should return a non-default color
		if color == styles.Text {
			t.Errorf("expected specific color for %s, got default", m)
		}
	}
}

func TestMethodColor_UnknownMethod(t *testing.T) {
	color := styles.MethodColor("TRACE")
	if color != styles.Text {
		t.Error("expected default color for unknown method")
	}
}

func TestMethod_ReturnsNonEmpty(t *testing.T) {
	result := styles.Method("GET")
	if result == "" {
		t.Error("expected non-empty styled string")
	}
}
