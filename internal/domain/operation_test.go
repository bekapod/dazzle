package domain_test

import (
	"testing"

	"dazzle/internal/domain"
)

func TestHTTPMethodConstants(t *testing.T) {
	tests := []struct {
		method domain.HTTPMethod
		want   string
	}{
		{domain.GET, "GET"},
		{domain.POST, "POST"},
		{domain.PUT, "PUT"},
		{domain.PATCH, "PATCH"},
		{domain.DELETE, "DELETE"},
		{domain.HEAD, "HEAD"},
		{domain.OPTIONS, "OPTIONS"},
	}

	for _, tt := range tests {
		if string(tt.method) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, tt.method)
		}
	}
}

func TestOperationFilter_ZeroValue(t *testing.T) {
	var filter domain.OperationFilter

	if filter.Query != "" {
		t.Error("expected empty query")
	}
	if filter.Method != "" {
		t.Error("expected empty method")
	}
	if filter.Tags != nil {
		t.Error("expected nil tags")
	}
}
