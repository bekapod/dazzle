package domain

import "testing"

func TestHTTPMethod_String(t *testing.T) {
	tests := []struct {
		method   HTTPMethod
		expected string
	}{
		{GET, "GET"},
		{POST, "POST"},
		{PUT, "PUT"},
		{PATCH, "PATCH"},
		{DELETE, "DELETE"},
	}

	for _, test := range tests {
		if string(test.method) != test.expected {
			t.Errorf("expected %s, got %s", test.expected, string(test.method))
		}
	}
}

func TestOperation_Creation(t *testing.T) {
	op := Operation{
		Path:    "/api/users",
		Method:  GET,
		Summary: "List all users",
		Tags:    []string{"users", "api"},
	}

	if op.Path != "/api/users" {
		t.Errorf("expected path '/api/users', got %s", op.Path)
	}

	if op.Method != GET {
		t.Errorf("expected method GET, got %s", op.Method)
	}

	if op.Summary != "List all users" {
		t.Errorf("expected summary 'List all users', got %s", op.Summary)
	}

	if len(op.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(op.Tags))
	}
}