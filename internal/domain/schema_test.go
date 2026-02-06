package domain_test

import (
	"testing"

	"dazzle/internal/domain"
)

func TestSchemaTypeConstants(t *testing.T) {
	tests := []struct {
		schemaType domain.SchemaType
		want       string
	}{
		{domain.SchemaTypeString, "string"},
		{domain.SchemaTypeNumber, "number"},
		{domain.SchemaTypeInteger, "integer"},
		{domain.SchemaTypeBoolean, "boolean"},
		{domain.SchemaTypeArray, "array"},
		{domain.SchemaTypeObject, "object"},
	}

	for _, tt := range tests {
		if string(tt.schemaType) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, tt.schemaType)
		}
	}
}

func TestParameterInConstants(t *testing.T) {
	tests := []struct {
		paramIn domain.ParameterIn
		want    string
	}{
		{domain.ParameterInPath, "path"},
		{domain.ParameterInQuery, "query"},
		{domain.ParameterInHeader, "header"},
		{domain.ParameterInCookie, "cookie"},
	}

	for _, tt := range tests {
		if string(tt.paramIn) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, tt.paramIn)
		}
	}
}
