package openapi_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"dazzle/internal/domain"
	"dazzle/internal/infrastructure/openapi"
)

func fixturesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "testdata", "fixtures")
}

func TestRepository_Load(t *testing.T) {
	repo := openapi.NewRepository()
	spec, err := repo.Load(context.Background(), filepath.Join(fixturesDir(), "petstore.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("spec info", func(t *testing.T) {
		if spec.Info.Title != "Petstore API" {
			t.Errorf("expected title 'Petstore API', got %q", spec.Info.Title)
		}
		if spec.Info.Version != "1.0.0" {
			t.Errorf("expected version '1.0.0', got %q", spec.Info.Version)
		}
	})

	t.Run("servers", func(t *testing.T) {
		if len(spec.Servers) != 1 {
			t.Fatalf("expected 1 server, got %d", len(spec.Servers))
		}
		if spec.Servers[0].URL != "https://api.petstore.example" {
			t.Errorf("unexpected server URL: %s", spec.Servers[0].URL)
		}
	})

	t.Run("operations count", func(t *testing.T) {
		if len(spec.Operations) != 4 {
			t.Errorf("expected 4 operations, got %d", len(spec.Operations))
		}
	})

	t.Run("operation fields", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		listPets, ok := ops["listPets"]
		if !ok {
			t.Fatal("missing listPets operation")
		}
		if listPets.Method != domain.GET {
			t.Errorf("expected GET, got %s", listPets.Method)
		}
		if listPets.Path != "/pets" {
			t.Errorf("expected /pets, got %s", listPets.Path)
		}
		if listPets.Summary != "List all pets" {
			t.Errorf("unexpected summary: %s", listPets.Summary)
		}
	})

	t.Run("operation-level parameters", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		listPets := ops["listPets"]
		if len(listPets.Parameters) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(listPets.Parameters))
		}
		param := listPets.Parameters[0]
		if param.Name != "limit" {
			t.Errorf("expected param name 'limit', got %q", param.Name)
		}
		if param.In != domain.ParameterInQuery {
			t.Errorf("expected param in 'query', got %q", param.In)
		}
	})

	t.Run("path-level parameters merged into operations", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		// getPet has no operation-level params but should inherit petId from path level
		getPet := ops["getPet"]
		if len(getPet.Parameters) != 1 {
			t.Fatalf("expected 1 parameter (from path level), got %d", len(getPet.Parameters))
		}
		if getPet.Parameters[0].Name != "petId" {
			t.Errorf("expected param name 'petId', got %q", getPet.Parameters[0].Name)
		}
		if getPet.Parameters[0].In != domain.ParameterInPath {
			t.Errorf("expected param in 'path', got %q", getPet.Parameters[0].In)
		}

		// deletePet should also inherit petId from path level
		deletePet := ops["deletePet"]
		if len(deletePet.Parameters) != 1 {
			t.Fatalf("expected 1 parameter (from path level), got %d", len(deletePet.Parameters))
		}
		if deletePet.Parameters[0].Name != "petId" {
			t.Errorf("expected param name 'petId', got %q", deletePet.Parameters[0].Name)
		}
	})
}

func TestRepository_Load_InvalidFile(t *testing.T) {
	repo := openapi.NewRepository()
	_, err := repo.Load(context.Background(), "nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func indexByID(ops []domain.Operation) map[string]domain.Operation {
	m := make(map[string]domain.Operation, len(ops))
	for _, op := range ops {
		m[op.ID] = op
	}
	return m
}
